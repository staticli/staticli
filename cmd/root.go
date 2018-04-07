package cmd

import (
	"github.com/wheresalice/cali"
)

var (
	// Define this here, then all other files in cmd can add subcommands to it
	cli = cali.NewCli("staticli")
)

func init() {
	cli.SetShort("Static site generator tooling")
	cli.SetLong("Static site generator tooling")
	cli.Flags().StringP("port", "p", "2000", "Which port should we expose on the host?")
	cli.BindFlags()
}

func Execute() {
	cli.Start()
}
