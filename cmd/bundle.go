package cmd

import "github.com/wheresalice/cali"
import "os/user"
import log "github.com/Sirupsen/logrus"
import (
	_ "github.com/pkg/errors"
)

func init() {

	command := cli.NewCommand("bundle")
	command.SetShort("Run bundler")
	command.SetLong(`Starts a container and runs bundler.
Examples:
  To update gems.
  # staticli bundle update
Any addtional flags sent to the rake command come after the --, e.g.
  # staticli bundle install -- --path=_vendor
`)

	image := "staticli/rake"
	task := command.Task(image)
	command.Flags().StringP("tag", "t", "latest", "Tag (Ruby version) to use (ruby2.4, ruby2.5)")
	command.BindFlags()
	task.Conf.Entrypoint = []string{"bundle"}
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}
	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	task.SetInitFunc(func(t *cali.Task, args []string) {
		t.SetImage(image + ":" + cli.FlagValues().GetString("tag"))
		log.Infof("Using Ruby %s", cli.FlagValues().GetString("tag"))
	})
}
