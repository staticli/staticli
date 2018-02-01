package cmd

import "github.com/skybet/cali"
import "os/user"
import log "github.com/Sirupsen/logrus"
import "path"
import (
	_ "github.com/pkg/errors")

func init() {

	command := cli.NewCommand("heroku")
	command.SetShort("Run Heroku CLI tools")
	command.SetLong(`Starts a container and runs Heroku CLI tooling'.
Examples:
  # staticli heroku
`)
	command.BindFlags()

	task := command.Task("wingrunr21/alpine-heroku-cli")

	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}

	netRc := path.Join(u.HomeDir, ".netrc")

	netRcBind, err := task.Bind(netRc, "/root/.netrc")
	if err != nil {
		log.Fatalf("Unable to bind ~/.netrc: %s", err)
	}
	task.AddBind(netRcBind)

	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	task.SetInitFunc(func(t *cali.Task, args []string) {
	})
}
