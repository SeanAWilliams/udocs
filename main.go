package main

import "github.com/UltimateSoftware/udocs/cli/cmd"

//go:generate vfsgendev -source="github.com/UltimateSoftware/udocs/static".Assets

// buildNumber is set via -ldflags
var buildNumber string

func main() {
	// static.FS = http.Dir("static")
	// static.Generate()
	cmd.BuildNumber = buildNumber

	// commands MUST be in alphabetical order
	cmd.Root.AddCommand(
		cmd.Build(),
		cmd.Destroy(),
		cmd.Env(),
		cmd.Publish(),
		cmd.Serve(),
		cmd.Tar(),
		cmd.Validate(),
		cmd.Version(),
	)
	cmd.Root.SetHelpTemplate(helpTmpl)
	cmd.Root.Execute()
}

var helpTmpl = `
Description:
  {{with or .Long .Short }}{{. | trim}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
