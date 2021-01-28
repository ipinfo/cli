package main

import (
	"strings"
)

var appHelpMsg = strings.TrimLeft(`
Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[<global-opts>]{{end}}{{if .Commands}} <cmd> [<opts>]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[<args>]{{end}}{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description | nindent 3 | trim}}{{end}}

<cmd>:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

<global-opts>:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}
`, "\n")

var cmdHelpMsg = strings.TrimLeft(`
Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [<opts>]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[<args>]{{end}}{{end}}{{if .Description}}

Description:
   {{.Description | nindent 3 | trim}}{{end}}{{if .VisibleFlags}}

<opts>:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}
`, "\n")

var subcmdHelpMsg = strings.TrimLeft(`
Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} <cmd>{{if .VisibleFlags}} [<opts>]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[<args>]{{end}}{{end}}{{if .Description}}

Description:
   {{.Description | nindent 3 | trim}}{{end}}

<cmd>:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

<opts>:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}
`, "\n")
