package cmd

import (
	"github.com/skybet/cali"
	"github.com/docker/go-connections/nat"
	log "github.com/Sirupsen/logrus"
	_ "github.com/pkg/errors"
)

func init() {
	command := cli.NewCommand("skipper")
	command.SetShort("Run skipper")
	command.SetLong(`Run Skipper, an http router and reverse proxy.
Usage:
# staticli skipper -- -routes-file example.eskip
`)

	task := command.Task("registry.opensource.zalan.do/pathfinder/skipper:latest")
	task.Conf.Entrypoint = []string{"skipper"}
	task.SetInitFunc(func(t *cali.Task, args []string) {
		log.Infof("Serving http on port %s - http://127.0.0.1:%s", cli.FlagValues().GetString("port"), cli.FlagValues().GetString("port"))

		task.HostConf.PortBindings = nat.PortMap{
			nat.Port("9090/tcp"): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: cli.FlagValues().GetString("port")},
			},
		}

	})
}
