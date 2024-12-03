package lws

import (
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:    "lws",
	Aliases: []string{"letter-writing-set"},
	Usage:   "Manage letter writing set",
	Subcommands: []*cli.Command{
		createCommand,
		statusCommand,
	},
}
