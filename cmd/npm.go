package cmd

import "github.com/skybet/cali"
import "os/user"
import log "github.com/Sirupsen/logrus"

func init() {

	command := cli.NewCommand("npm")
	command.SetShort("Run npm")
	command.SetLong(`Starts a container and runs npm.
Examples:
  # staticli npm install
Any addtional flags sent to the npm command come after the --, e.g.
  # staticli npm <command> -- --key value
`)

	task := command.Task("node:alpine")
	task.Conf.Entrypoint = []string{"npm"}
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}
	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	task.SetInitFunc(func(t *cali.Task, args []string) {
	})
}
