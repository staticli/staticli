package cmd

import "github.com/skybet/cali"
import "os/user"
import "os"
import log "github.com/Sirupsen/logrus"
import (
_ "github.com/pkg/errors")

func init() {

	command := cli.NewCommand("github-release")
	command.SetShort("Run github-release")
	command.SetLong(`Starts a container and runs github-release
This is primarily an internal thing to staticli so that it can release itself, however there's no reason you can't use it yourself
Examples:
# staticli github-release "v1.0" *.tar.gz -- --github-access-token [...]
you can also export GITHUB_RELEASE_ACCESS_TOKEN
`)
	command.BindFlags()

	task := command.Task("buildkite/github-release")

	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}
	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	if os.Getenv("GITHUB_RELEASE_ACCESS_TOKEN")!=""{
		task.AddEnv("GITHUB_RELEASE_ACCESS_TOKEN", os.Getenv("GITHUB_RELEASE_ACCESS_TOKEN"))
	}
	task.SetInitFunc(func(t *cali.Task, args []string) {

	})
}
