package cmd

import "github.com/skybet/cali"
import "os/user"
import log "github.com/Sirupsen/logrus"

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
	})
}
