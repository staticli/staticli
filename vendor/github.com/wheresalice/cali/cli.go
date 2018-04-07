package cali

import (
	"fmt"
	"os"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
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

// cobraFunc represents the function signiture which cobra uses for it's Run, PreRun, PostRun etc.
type cobraFunc func(cmd *cobra.Command, args []string)

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
