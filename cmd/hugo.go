package cmd

import "github.com/skybet/cali"
import "os/user"
import log "github.com/Sirupsen/logrus"
import (
	_ "github.com/pkg/errors"
	"github.com/docker/go-connections/nat"
)

func init() {

	command := cli.NewCommand("hugo")
	command.SetShort("Run hugo")
	command.SetLong(`Starts a container and runs hugo.
Examples:
  To create a new site.
  # staticli hugo new site my-site
Any addtional flags sent to the hugo command come after the --, e.g.
  # staticli hugo serve -- --bind=0.0.0.0
`)

	task := command.Task("jojomi/hugo")
	task.Conf.Entrypoint = []string{"hugo"}
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}
	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	task.SetInitFunc(func(t *cali.Task, args []string) {
		log.Infof("Serving http on port %s - http://127.0.0.1:%s", cli.FlagValues().GetString("port"), cli.FlagValues().GetString("port"))
		task.HostConf.PortBindings = nat.PortMap{
			nat.Port("1313/tcp"): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: cli.FlagValues().GetString("port")},
			},
		}

	})
}
