package cmd

import "github.com/skybet/cali"
import "os/user"
import log "github.com/Sirupsen/logrus"

func init() {

	command := cli.NewCommand("surge [command]")
	command.SetShort("Run surge.sh")
	command.SetLong(`Starts a container and runs surge.sh.
Examples:
  To publish the current directory.
  # staticli surge

  You can also run subcommands.
  # staticli surge list
`)

	task := command.Task("staticli/surge")
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}
	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	task.SetInitFunc(func(t *cali.Task, args []string) {

	})
}
