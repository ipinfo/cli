package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	app.Version = "0.1.0"
	app.EnableBashCompletion = true
	app.CommandNotFound = func(c *cli.Context, cmd string) {
		fmt.Printf("err: \"%s\" is not a command.\n", cmd)
		fmt.Println()
		cli.ShowAppHelp(c)
	}
	app.HideHelpCommand = true
	app.UseShortOptionHandling = true
	app.Action = func(c *cli.Context) error {
		// If there are no arguments, print normal help text.
		args := c.Args()
		if !args.Present() {
			cli.ShowAppHelp(c)
			return nil
		}

		// Check for whether the "command" is really an IP or ASN; handle those
		// properly.
		if err := cmdIP(c); err != errNotIP {
			return err
		}
		if err := cmdAsn(c); err != errNotASN {
			return err
		}

		return cli.ShowCommandHelp(c, args.First())
	}
	app.Commands = []*cli.Command{
		{
			Name:  "myip",
			Usage: "get details for your IP",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "json",
					Aliases: []string{"j"},
					Usage:   "output JSON format",
				},
			},
			Action: cmdMyIP,
		},
		{
			Name:  "completion",
			Usage: "generate auto-completion script for a shell environment",
			Subcommands: []*cli.Command{
				{
					Name:   "bash",
					Usage:  "generate bash auto-completion",
					Action: cmdCompletionBash,
				},
				{
					Name:   "zsh",
					Usage:  "generate ZSH auto-completion",
					Action: cmdCompletionZsh,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
