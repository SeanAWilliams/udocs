package main

import "github.com/gdscheele/udocs/cli/cmd"

var version string // set via -ldflags

func main() {
	cmd.VersionNumber = version

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
