package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/wheresalice/cali"

	"github.com/staticli/staticli/lib"
)

func init() {

	command := cli.NewCommand("version")
	command.SetShort("Which version are we running?")

	var taskFunc cali.TaskFunc = func(t *cali.Task, args []string) {
		lib.PrintVersion()

		isLatestVersion, releaseData, err := lib.IsLatestVersion()
		if err != nil {
			log.Fatalf("Unable to check for update: %s", err)
		}

		if !isLatestVersion {
			log.Infof("You're not running the latest version 😱")
			log.Infof("Update to v%s with: staticli update", releaseData.Name)
		}
	}

	// Simple task, just runs a function
	command.Task(taskFunc)

}
