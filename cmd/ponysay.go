package cmd

import "github.com/skybet/cali"

func init() {

	command := cli.NewCommand("ponysay")
	command.SetShort("Run ponysay")
	command.SetLong(`Starts a container and runs ponysay.
Examples:
  # staticli ponysay "hello world"
Any additional flags sent to the ponysay command come after the --, e.g.
  # staticli ponysay -- -q
`)

	task := command.Task("mpepping/ponysay")
	task.Conf.Entrypoint = []string{"ponysay"}
	task.SetInitFunc(func(t *cali.Task, args []string) {
	})
}
