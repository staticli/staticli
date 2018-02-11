package cmd

import "github.com/skybet/cali"
import "os/user"
import log "github.com/Sirupsen/logrus"
import (
	_ "github.com/pkg/errors"
)

func init() {

	command := cli.NewCommand("gulp")
	command.SetShort("Run gulp")
	command.SetLong(`Starts a container and runs gulp
Examples:
  To run gulp watch:
  # staticli gulp
  Any addtional flags sent to the rake command come after the --, e.g.
  # staticli rake preview -- --future
`)
	command.Flags().StringP("gulp_task", "t", "watch", "Gulp task to run")
	command.BindFlags()

	task := command.Task("agomezmoron/docker-gulp")

	src, err := task.Bind(".", "/src")
	if err != nil {
		log.Fatalf("Unable to bind . : %s", err)
	}
	task.AddBind(src)

	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}
	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	task.AddEnv("GULP_TASK", cli.FlagValues().GetString("gulp_task"))
	task.SetInitFunc(func(t *cali.Task, args []string) {

	})
}
