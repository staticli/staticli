package cmd

import "github.com/skybet/cali"
import "github.com/docker/go-connections/nat"
import "os/user"
import log "github.com/Sirupsen/logrus"
import (
	_ "github.com/pkg/errors")

func init() {

	command := cli.NewCommand("simplehttp")
	command.SetShort("Run Python SimpleHTTPServer")
	command.SetLong(`Starts a container and runs Python's SimpleHTTPServer module to serve the current directory over http'.
Examples:
  To serve the current directory on http://127.0.0.1:8000.
  # staticli simplehttp
  You can also set the port to listen on if you want something other than port 8000.
  # staticli simplehttp -p 3000
`)
	command.Flags().StringP("port", "p", "8000", "Port to expose on host")
	command.BindFlags()

	task := command.Task("staticli/simplehttp")
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to find uid for user: %s", err)
	}
	task.AddEnv("HOST_USER_ID", u.Uid)
	task.AddEnv("HOST_GROUP_ID", u.Gid)
	task.SetInitFunc(func(t *cali.Task, args []string) {

		task.HostConf.PortBindings = nat.PortMap{
			nat.Port("8000/tcp"): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: cli.FlagValues().GetString("port")},
			},
		}

	})
}
