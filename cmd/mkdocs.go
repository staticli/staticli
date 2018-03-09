package cmd

import (
	"github.com/skybet/cali"
	"github.com/docker/go-connections/nat"
	log "github.com/Sirupsen/logrus"
	_ "github.com/pkg/errors"
)

func init() {
	command := cli.NewCommand("mkdocs")
	command.SetShort("Run MkDocs")
	command.SetLong(`Run MkDocs, a fast and simple static site generator that's geared towards building project documentation'.
Usage:
# staticli mkdocs serve
# staticli mkdocs new <directory_name>
# staticli mkdocs build
`)

	task := command.Task("moird/mkdocs")
	task.Conf.Entrypoint = []string{"mkdocs"}
	task.SetInitFunc(func(t *cali.Task, args []string) {
		log.Infof("Serving http on port %s - http://127.0.0.1:%s", cli.FlagValues().GetString("port"), cli.FlagValues().GetString("port"))
		log.Warn("This will only work if you have set 'dev_addr: 0.0.0.0:8000' in mkdocs.yml")

		task.HostConf.PortBindings = nat.PortMap{
			nat.Port("8000/tcp"): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: cli.FlagValues().GetString("port")},
			},
		}

	})
}
