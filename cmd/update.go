package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/wheresalice/cali"
	"github.com/staticli/staticli/lib"
)

func init() {

	command := cli.NewCommand("update")
	command.SetShort("Update the current running version of lucli")
	command.Flags().BoolP("images","i",false, "Should we also update local images?")
	command.BindFlags()

	var taskFunc cali.TaskFunc = func(t *cali.Task, args []string) {

		if cli.FlagValues().GetBool("images") {
			docker := cali.NewDockerClient()

			images := []string{"ketouem/ag-alpine", "staticli/rake", "buildkite/github-release", "agomezmoron/docker-gulp", "wingrunr21/alpine-heroku-cli", "jojomi/hugo", "jekyll/jekyll", "moird/mkdocs", "node:alpine", "mpepping/ponysay", "staticli/proselint", "staticli/simplehttp", "registry.opensource.zalan.do/pathfinder/skipper:latest", "staticli/surge"}
			for _, image := range images {
				imageExists, _ := docker.ImageExists(image)
				if imageExists {
					log.Debugf("checking for updates to %s", image)
					docker.PullImage(image)
				} else {
					log.Debugf("%s image doesn't exist locally, not pulling", image)
				}

			}
		}

		lib.PrintVersion()

		isLatestVersion, updateReleaseData, err := lib.IsLatestVersion()
		if err != nil {
			log.Fatalf("Unable to check for update: %s", err)
		}

		if !isLatestVersion {
			err = lib.Update(updateReleaseData)
			if err != nil {
				log.Fatalf("Unable to update: %s", err)
			}
		}
	}

	command.Task(taskFunc)

}