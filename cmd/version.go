package cmd

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"encoding/json"
	"github.com/skybet/cali"
	"github.com/staticli/staticli/lib"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

type Release struct {
	Name string
}

func init() {

	command := cli.NewCommand("version")
	command.SetShort("Show current staticli version")

	releaseUrl := "https://api.github.com/repos/staticli/staticli/releases/latest"

	var taskFunc cali.TaskFunc = func(t *cali.Task, args []string) {
		log.Infof("staticli v%s (%s) (%s)", lib.Version, lib.BuildTime, lib.BuildCommit)

		releaseData := Release{}
		getJson(releaseUrl, &releaseData)
		if releaseData.Name != "" {
			log.Infof("staticli v%s is the latest release", releaseData.Name)
			if lib.Version != releaseData.Name {
				log.Warnf("You should update your version of staticli")
			}

		}
	}

	command.Task(taskFunc)

}
