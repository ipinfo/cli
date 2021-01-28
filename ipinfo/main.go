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

		// check whether initial command is an IP or ASN.
		ipOrASN := args.First()
		if isIP(ipOrASN) {
			ipStr := ipOrASN
			ipCmd := c.App.Command("_ip")
			ipCmd.Name = ipStr
			ipCmd.HelpName = progBase + " " + ipStr

			newArgs := []string{os.Args[0], ipStr, "--ip", ipStr}
			newArgs = append(newArgs, os.Args[2:]...)
			return c.App.Run(newArgs)
		}
		if isASN(ipOrASN) {
			asn := ipOrASN
			ipCmd := c.App.Command("_asn")
			ipCmd.Name = asn
			ipCmd.HelpName = progBase + " " + asn

			newArgs := []string{os.Args[0], asn, "--asn", asn}
			newArgs = append(newArgs, os.Args[2:]...)
			return c.App.Run(newArgs)
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
		/* hidden commands as hacks to allow ip/asn positional arguments
		   without requiring them to be behind commands that the user has to
		   input manually. */
		{
			Name:   "_ip",
			Hidden: true,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "ip",
					Required: true,
					Hidden:   true,
				},
				&cli.BoolFlag{
					Name:    "json",
					Aliases: []string{"j"},
					Usage:   "output JSON format",
				},
			},
			Action: cmdIP,
		},
		{
			Name:   "_asn",
			Hidden: true,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "asn",
					Required: true,
					Hidden:   true,
				},
				&cli.BoolFlag{
					Name:    "json",
					Aliases: []string{"j"},
					Usage:   "output JSON format",
				},
			},
			Action: cmdASN,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
