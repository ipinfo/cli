package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/urfave/cli/v2"
)

var progBase = filepath.Base(os.Args[0])

var ii *ipinfo.Client

func prepareIpinfoClient(c *cli.Context) error {
	tok := c.String("token")
	if tok == "" {
		tok, _ = restoreToken()
	}

	ii = ipinfo.NewClient(nil, nil, tok)
	return nil
}

func main() {
	tokenFlag := &cli.StringFlag{
		Name:    "token",
		Aliases: []string{"t"},
		Usage:   "use `TOK` as API token",
	}
	jsonFlag := &cli.BoolFlag{
		Name:    "json",
		Aliases: []string{"j"},
		Usage:   "output JSON format",
	}
	csvFlag := &cli.BoolFlag{
		Name:    "csv",
		Aliases: []string{"c"},
		Usage:   "output CSV format",
	}

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
		/* HACK
		the tactic used here uses hidden commands allow ip/asn positional
		arguments without requiring them to be behind commands that the user
		has to input manually.
		*/
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
			asn := strings.ToUpper(ipOrASN)
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
				jsonFlag,
				tokenFlag,
			},
			Before:          prepareIpinfoClient,
			Action:          cmdMyIP,
			HideHelpCommand: true,
		},
		{
			Name:      "bulk",
			Usage:     "get details for multiple IPs in bulk",
			ArgsUsage: "<paths or '-' or cidrs or ip-range>",
			Description: "Accepts file paths, '-' for stdin, CIDRs and IP ranges.\n" +
				"\n" +
				"# Lookup all IPs from stdin.\n" +
				"$ " + progBase + " bulk -\n" +
				"\n" +
				"# Lookup all IPs in 2 files.\n" +
				"$ " + progBase + " bulk /path/to/iplist1.txt /path/to/iplist2.txt\n" +
				"\n" +
				"# Lookup all IPs from CIDR.\n" +
				"$ " + progBase + " bulk 8.8.8.0/24\n" +
				"\n" +
				"# Lookup all IPs from multiple CIDRs.\n" +
				"$ " + progBase + " bulk 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24\n" +
				"\n" +
				"# Lookup all IPs in an IP range.\n" +
				"$ " + progBase + " bulk 8.8.8.0 8.8.8.255",
			Flags: []cli.Flag{
				jsonFlag,
				csvFlag,
			},
			Before:          prepareIpinfoClient,
			Action:          cmdBulk,
			HideHelpCommand: true,
		},
		{
			Name:            "prips",
			Usage:           "print IP list from CIDR or range",
			Action:          cmdPrips,
			HideHelpCommand: true,
			ArgsUsage:       "<cidrs or ip-range>",
			Description: "Accepts CIDRs (e.g. 8.8.8.0/24) or an IP range (e.g. 8.8.8.0 8.8.8.255).\n" +
				"\n" +
				"# List all IPs in a CIDR.\n" +
				"$ " + progBase + " prips 8.8.8.0/24\n" +
				"\n" +
				"# List all IPs in multiple CIDRs.\n" +
				"$ " + progBase + " prips 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24\n" +
				"\n" +
				"# List all IPs in an IP range.\n" +
				"$ " + progBase + " prips 8.8.8.0 8.8.8.255",
		},
		{
			Name:  "login",
			Usage: "save an API token session",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "token",
					Aliases: []string{"t"},
					Usage: "token to login with " +
						"(potentially unsafe; let CLI prompt you instead)",
				},
			},
			Action:          cmdLogin,
			HideHelpCommand: true,
		},
		{
			Name:            "logout",
			Usage:           "delete your current API token session",
			Action:          cmdLogout,
			HideHelpCommand: true,
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
			HideHelpCommand: true,
		},
		{
			Name:   "_ip",
			Hidden: true,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "ip",
					Required: true,
					Hidden:   true,
				},
				jsonFlag,
				tokenFlag,
			},
			Before: prepareIpinfoClient,
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
				tokenFlag,
			},
			Before: prepareIpinfoClient,
			Action: cmdASN,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
