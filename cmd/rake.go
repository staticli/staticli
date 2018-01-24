package cmd

import "github.com/skybet/cali"

func init() {

	command := cli.NewCommand("rake")
	command.SetShort("Run rake preview")
	command.SetLong(`Starts a container and runs rake preview.
Examples:
  To render the site locally.
  # staticli rake
  Any addtional flags sent to the rake command come after the --, e.g.
  # staticli rake preview -- --future
`)
	command.BindFlags()

	task := command.Task("kaerast/rake-preview")
	task.HostConf.PublishAllPorts = true
	// @todo the above publishes ports defined in the Dockerfile as `EXPOSE` on a random port.  We would prefer a defined port, however we don't understand enough Go to use task.HostConf.PortBindings
	task.SetInitFunc(func(t *cali.Task, args []string) {})
}
