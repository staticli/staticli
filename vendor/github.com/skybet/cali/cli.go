package cali

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	EXIT_CODE_RUNTIME_ERROR = 1
	EXIT_CODE_API_ERROR     = 2

	workdir = "/tmp/workspace"
)

var (
	debug, jsonLogs, nonInteractive bool
	dockerHost                      string
	myFlags                         *viper.Viper
	gitCfg                          *GitCheckoutConfig
)

// TaskFunc is a function executed by a Task when the command the Task belongs to is run
type TaskFunc func(t *Task, args []string)

// defaultTaskFunc is the TaskFunc which is executed unless a custom TaskFunc is
// attached to the Task
var defaultTaskFunc TaskFunc = func(t *Task, args []string) {
	if err := t.SetDefaults(args); err != nil {
		log.Fatalf("Error setting container defaults: %s", err)
	}
	if err := t.InitDocker(); err != nil {
		log.Fatalf("Error initialising Docker: %s", err)
	}
	if _, err := t.StartContainer(true, ""); err != nil {
		log.Fatalf("Error executing task: %s", err)
	}
}

// Task is the action performed when it's parent command is run
type Task struct {
	f, init TaskFunc
	*DockerClient
}

// SetFunc sets the TaskFunc which is run when the parent command is run
// if this is left unset, the defaultTaskFunc will be executed instead
func (t *Task) SetFunc(f TaskFunc) {
	t.f = f
}

// SetInitFunc sets the TaskFunc which is executed before the main TaskFunc. It's
// pupose is to do any setup of the DockerClient which depends on command line args
// for example
func (t *Task) SetInitFunc(f TaskFunc) {
	t.init = f
}

// SetDefaults sets the default host config for a task container
// Mounts the PWD to /tmp/workspace
// Mounts your ~/.aws directory to /root - change this if your image runs as a non-root user
// Sets /tmp/workspace as the workdir
// Configures git
func (t *Task) SetDefaults(args []string) error {
	t.SetWorkDir(workdir)
	awsDir, err := t.Bind("~/.aws", "/root/.aws")
	if err != nil {
		return err
	}
	t.AddBinds([]string{awsDir})

	err = t.BindFromGit(gitCfg, func() error {
		pwd, err := t.Bind("./", workdir)
		if err != nil {
			return err
		}
		t.AddBinds([]string{pwd})
		return nil
	})
	if err != nil {
		return err
	}
	t.SetCmd(args)
	return nil
}

// Bind is a utility function which will return the correctly formatted string when given a source
// and destination directory
//
// The ~ symbol and relative paths will be correctly expanded depending on the host OS
func (t *Task) Bind(src, dst string) (string, error) {
	var expanded string

	if strings.HasPrefix(src, "~") {
		usr, err := user.Current()

		if err != nil {
			return expanded, fmt.Errorf("Error expanding bind path: %s")
		}
		expanded = filepath.Join(usr.HomeDir, src[2:])
	} else {
		expanded = src
	}
	expanded, err := filepath.Abs(expanded)

	if err != nil {
		return expanded, fmt.Errorf("Error expanding bind path: %s")
	}
	return fmt.Sprintf("%s:%s", expanded, dst), nil
}

// cobraFunc represents the function signiture which cobra uses for it's Run, PreRun, PostRun etc.
type cobraFunc func(cmd *cobra.Command, args []string)

// command is the actual command run by the cli and essentially just wraps cobra.Command and
// has an associated Task
type Command struct {
	name    string
	RunTask *Task
	cobra   *cobra.Command
}

// SetShort sets the short description of the command
func (c *Command) SetShort(s string) {
	c.cobra.Short = s
}

// SetLong sets the long description of the command
func (c *Command) SetLong(l string) {
	c.cobra.Long = l
}

// setPreRun sets the cobra.Command.PreRun function
func (c *Command) setPreRun(f cobraFunc) {
	c.cobra.PreRun = f
}

// setRun sets the cobra.Command.Run function
func (c *Command) setRun(f cobraFunc) {
	c.cobra.Run = f
}

// Task is something executed by a command
func (c *Command) Task(def interface{}) *Task {
	t := &Task{DockerClient: NewDockerClient()}

	switch d := def.(type) {
	case string:
		t.SetImage(d)
		t.SetFunc(defaultTaskFunc)
	case TaskFunc:
		t.SetFunc(d)
	default:
		// Slightly unidiomatic to blow up here rather than return an error
		// choosing to so as to keep the API uncluttered and also if you get here it's
		// an implementation error rather than a runtime error.
		fmt.Println("Unknown Task type. Must either be an image (string) or a TaskFunc")
		os.Exit(EXIT_CODE_API_ERROR)
	}
	c.RunTask = t
	return t
}

// Flags returns the FlagSet for the command and is used to set new flags for the command
func (c *Command) Flags() *flag.FlagSet {
	return c.cobra.PersistentFlags()
}

// BindFlags needs to be called after all flags for a command have been defined
func (c *Command) BindFlags() {
	c.Flags().VisitAll(func(f *flag.Flag) {
		myFlags.BindPFlag(f.Name, f)
		myFlags.SetDefault(f.Name, f.DefValue)
	})
}

// commands is a set of commands
type commands map[string]*Command

// Cli is the application itself
type Cli struct {
	name    string
	cfgFile *string
	cmds    commands
	*Command
}

// NewCli returns a brand new cli
func NewCli(n string) *Cli {
	c := Cli{
		name: n,
		cmds: make(commands),
		Command: &Command{
			name:  n,
			cobra: &cobra.Command{Use: n},
		},
	}
	c.cobra.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if jsonLogs {
			log.SetFormatter(&log.JSONFormatter{})
		}
	}
	myFlags = viper.New()
	return &c
}

// NewCommand returns a brand new command attached to it's parent cli
func (c *Cli) NewCommand(n string) *Command {
	cmd := &Command{
		name:  n,
		cobra: &cobra.Command{Use: n},
	}
	c.cmds[n] = cmd

	cmd.setPreRun(func(c *cobra.Command, args []string) {
		// PreRun function is optional
		if cmd.RunTask.init != nil {
			cmd.RunTask.init(cmd.RunTask, args)
		}
	})
	cmd.setRun(func(c *cobra.Command, args []string) {
		cmd.RunTask.f(cmd.RunTask, args)
	})
	c.cobra.AddCommand(cmd.cobra)
	return cmd
}

// FlagValues returns the wrapped viper object allowing the API consumer to use methods
// like GetString to get values from config
func (c *Cli) FlagValues() *viper.Viper {
	return myFlags
}

// initFlags does the intial setup of the root command's persistent flags
func (c *Cli) initFlags() {
	var cfg string
	txt := fmt.Sprintf("config file (default is $HOME/.%s.yaml)", c.name)
	c.cobra.PersistentFlags().StringVar(&cfg, "config", "", txt)
	c.cfgFile = &cfg

	var dockerSocket string
	if runtime.GOOS == "windows" {
		dockerSocket = "npipe:////./pipe/docker_engine"
	} else {
		dockerSocket = "unix:///var/run/docker.sock"
	}
	c.Flags().StringVarP(&dockerHost, "docker-host", "H", dockerSocket, "URI of Docker Daemon")
	myFlags.BindPFlag("docker-host", c.Flags().Lookup("docker-host"))
	myFlags.SetDefault("docker-host", dockerSocket)

	c.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	myFlags.BindPFlag("debug", c.Flags().Lookup("debug"))
	myFlags.SetDefault("debug", true)

	c.Flags().BoolVarP(&jsonLogs, "json", "j", false, "Log in json format")
	myFlags.BindPFlag("json", c.Flags().Lookup("json"))
	myFlags.SetDefault("json", true)

	c.Flags().BoolVarP(&nonInteractive, "non-interactive", "N", false, "Do not create a tty for Docker")
	myFlags.BindPFlag("non-interactive", c.Flags().Lookup("non-interactive"))
	myFlags.SetDefault("non-interactive", false)

	gitCfg = new(GitCheckoutConfig)
	c.Flags().StringVarP(&gitCfg.Repo, "git", "g", "", "Git repo to checkout and build. Default behaviour is to build $PWD.")
	myFlags.BindPFlag("git", c.Flags().Lookup("git"))

	c.Flags().StringVarP(&gitCfg.Branch, "git-branch", "b", "master", "Branch to checkout. Only makes sense when combined with the --git flag.")
	myFlags.BindPFlag("branch", c.Flags().Lookup("branch"))
	myFlags.SetDefault("branch", "master")

	c.Flags().StringVarP(&gitCfg.RelPath, "git-path", "P", "", "Path within a git repo where we want to operate.")
	myFlags.BindPFlag("git-path", c.Flags().Lookup("git-path"))
}

// initConfig does the initial setup of viper
func (c *Cli) initConfig() {
	if *c.cfgFile != "" {
		myFlags.SetConfigFile(*c.cfgFile)
	} else {
		myFlags.SetConfigName(fmt.Sprintf(".%s", c.name))
		myFlags.AddConfigPath(".")     // First check current working directory
		myFlags.AddConfigPath("$HOME") // Fallback to home directory, if that is not set
	}
	myFlags.AutomaticEnv()

	// If a config file is found, read it in.
	myFlags.ReadInConfig()
	// Above returns an error if it doesn't find a config file
	// But we don't care
}

// Start the fans please!
func (c *Cli) Start() {
	c.initFlags()
	cobra.OnInitialize(c.initConfig)

	if err := c.cobra.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(EXIT_CODE_RUNTIME_ERROR)
	}
}
