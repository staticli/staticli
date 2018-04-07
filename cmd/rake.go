package cmd

import "github.com/wheresalice/cali"
import "github.com/docker/go-connections/nat"
import "os/user"
import log "github.com/Sirupsen/logrus"
import (
	_ "github.com/pkg/errors"
)

func init() {

	command := cli.NewCommand("rake")
	command.SetShort("Run rake preview")
	command.SetLong(`Starts a container and runs rake preview.
Examples:
  To render the site locally.
  # staticli rake
  Any addtional flags sent to the rake command come after the --, e.g.
  # staticli rake preview -- --future
`)

	image := "staticli/rake"
	task := command.Task(image)
	command.Flags().StringP("tag", "t", "latest", "Tag (Ruby version) to use (latest, ruby2.4, ruby2.5)")
	command.BindFlags()
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}
	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	task.SetInitFunc(func(t *cali.Task, args []string) {
		t.SetImage(image + ":" + cli.FlagValues().GetString("tag"))
		log.Infof("Using Tag %s", cli.FlagValues().GetString("tag"))
		log.Infof("Serving http on port %s - http://127.0.0.1:%s", cli.FlagValues().GetString("port"), cli.FlagValues().GetString("port"))
		task.HostConf.PortBindings = nat.PortMap{
			nat.Port("4000/tcp"): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: cli.FlagValues().GetString("port")},
			},
		}

	})
}
