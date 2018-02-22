package cmd

import "github.com/skybet/cali"
import "github.com/docker/go-connections/nat"
import "os/user"
import log "github.com/Sirupsen/logrus"
import (
	_ "github.com/pkg/errors"
)

func init() {

	command := cli.NewCommand("jekyll")
	command.SetShort("Run jekyll")
	command.SetLong(`Starts a container and runs jekyll.
Examples:
  # staticli jekyll new my-awesome-site
Any addtional flags sent to the npm command come after the --, e.g.
  # staticli jekyll <command> -- --key value
`)

	task := command.Task("jekyll/jekyll")
	task.Conf.Entrypoint = []string{"jekyll"}
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}
	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	task.SetInitFunc(func(t *cali.Task, args []string) {
		log.Infof("Serving http on port %s - http://127.0.0.1:%s", cli.FlagValues().GetString("port"), cli.FlagValues().GetString("port"))
		task.HostConf.PortBindings = nat.PortMap{
			nat.Port("4000/tcp"): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: cli.FlagValues().GetString("port")},
			},
		}

	})
}
