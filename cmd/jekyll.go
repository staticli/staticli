package cmd

import "github.com/wheresalice/cali"
import "github.com/docker/go-connections/nat"
import "os/user"
import log "github.com/Sirupsen/logrus"
import (
	"io"
	"os"
	"path"

	_ "github.com/pkg/errors"
)

func IsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}

func init() {

	command := cli.NewCommand("jekyll")
	command.SetShort("Run jekyll")
	command.SetLong(`Starts a container and runs jekyll.
Examples:
  # staticli jekyll new my-awesome-site
Any additional flags sent to the jekyll command come after the --, e.g.
  # staticli jekyll <command> -- --key value
`)

	task := command.Task("jekyll/jekyll")
	task.Conf.Entrypoint = []string{"jekyll"}
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}

	task.AddEnv("JEKYLL_UID", u.Uid)
	task.AddEnv("JEKYLL_GID", u.Gid)

	cacheDir := path.Join(u.HomeDir, ".cache", "staticli", "jekyll", "bundle")
	os.Mkdir(cacheDir, 0700)

	cacheBind, err := task.Bind(cacheDir, "/usr/local/bundle")
	if err != nil {
		log.Fatalf("Unable to bind ~/.cache/staticli/ruby: %s", err)
	}
	task.AddBind(cacheBind)

	task.SetInitFunc(func(t *cali.Task, args []string) {
		os.Mkdir(cacheDir, 0700)
		if IsEmpty(cacheDir) {
			log.Warnf("You'll need to be online to run this because your cache directory(%s) is empty", cacheDir)
		}
		log.Debugf("Using %s for bundle cache directory", cacheDir)

		if len(args) >= 1 && args[0] == "build" {
			log.Debugf("Not mapping a port. Unrequired for a build")
		} else {

			log.Infof("Serving http on port %s - http://127.0.0.1:%s", cli.FlagValues().GetString("port"), cli.FlagValues().GetString("port"))

			task.HostConf.PortBindings = nat.PortMap{
				nat.Port("4000/tcp"): []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: cli.FlagValues().GetString("port")},
				},
			}
		}

	})
}
