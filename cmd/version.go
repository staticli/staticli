package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/skybet/cali"
	"github.com/staticli/staticli/lib"
	"runtime"
)

func init() {

	command := cli.NewCommand("version")
	command.SetShort("Show current staticli version")

	releaseUrl := "https://api.github.com/repos/staticli/staticli/releases/latest"

	var taskFunc cali.TaskFunc = func(t *cali.Task, args []string) {
		log.Infof("staticli.%s.%s v%s (%s) (%s)", runtime.GOOS, runtime.GOARCH, lib.Version, lib.BuildTime, lib.BuildCommit)

		releaseData := lib.Release{}
		lib.GetJson(releaseUrl, &releaseData)
		if releaseData.Name != "" {
			log.Infof("staticli v%s is the latest release", releaseData.Name)
			if lib.Version != releaseData.Name {
				log.Warnf("You should update your version of staticli")
			}

		}
	}

	command.Task(taskFunc)

}
