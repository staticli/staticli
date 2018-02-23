package cmd

import "github.com/skybet/cali"

func init() {

	command := cli.NewCommand("ag")
	command.SetShort("Run ag")
	command.SetLong(`Starts a container and runs Silver Surfer.
Examples:
  # staticli ag foo
Any addtional flags sent to the npm command come after the --, e.g.
  # staticli ag -- --help
  # staticli ag -- --markdown foo
`)

	task := command.Task("ketouem/ag-alpine")
	task.Conf.Entrypoint = []string{"/the_silver_searcher/ag"}
	task.SetInitFunc(func(t *cali.Task, args []string) {
	})
}
