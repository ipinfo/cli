package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var progBase = filepath.Base(os.Args[0])

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n", c.App.Version)
	}
	cli.AppHelpTemplate = appHelpMsg
	cli.CommandHelpTemplate = cmdHelpMsg
	cli.SubcommandHelpTemplate = subcmdHelpMsg

	app := cli.NewApp()
	app.Name = progBase
	app.Usage = "CLI for the IPinfo API"
	app.Version = "0.1.0"
	app.EnableBashCompletion = true
	app.HideHelpCommand = true
	app.UseShortOptionHandling = true
	app.Commands = []*cli.Command{
		{
			Name:  "completion",
			Usage: "generate auto-completion script for a shell environment",
			Subcommands: []*cli.Command{
				{
					Name:  "bash",
					Usage: "generate bash auto-completion",
					Action: func(c *cli.Context) error {
						fmt.Printf(
							strings.ReplaceAll(
								autocompleteBash, "$$PROG$$", progBase,
							),
						)
						return nil
					},
				},
				{
					Name:  "zsh",
					Usage: "generate ZSH auto-completion",
					Action: func(c *cli.Context) error {
						fmt.Printf(
							strings.ReplaceAll(
								autocompleteZsh, "$$PROG$$", progBase,
							),
						)
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
