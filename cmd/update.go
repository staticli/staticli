package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/skybet/cali"
	"github.com/staticli/staticli/lib"
	"runtime"
	"net/http"
	"github.com/inconshreveable/go-update"
	"fmt"
)

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		log.Errorf("Couldn't update staticli: %s", err)
	}
	return err
}

func init() {

	command := cli.NewCommand("update")
	command.SetShort("Update the current running version of staticli")

	releaseUrl := "https://api.github.com/repos/staticli/staticli/releases/latest"

	var taskFunc cali.TaskFunc = func(t *cali.Task, args []string) {
		log.Infof("you are running staticli v%s (%s) (%s)", lib.Version, lib.BuildTime, lib.BuildCommit)

		releaseData := lib.Release{}
		lib.GetJson(releaseUrl, &releaseData)
		if releaseData.Name == "" {
			log.Error("Couldn't get latest version of staticli, do you have an internet connection?")
		} else {
			log.Infof("staticli v%s is the latest release", releaseData.Name)
			if lib.Version != releaseData.Name {
				log.Infof("updating staticli to v%s", releaseData.Name)
				// @todo this url is subject to change so we should probably get it from the github json
				updateUrl := fmt.Sprintf("https://github.com/staticli/staticli/releases/download/%s/staticli.%s.%s", releaseData.Name, runtime.GOOS, runtime.GOARCH)
				doUpdate(updateUrl)
				log.Info("update complete")
			}

		}
	}

	command.Task(taskFunc)

}
