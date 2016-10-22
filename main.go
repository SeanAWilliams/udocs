package main

import (
	"github.com/ultimatesoftware/udocs/cli/cmd"
	"github.com/ultimatesoftware/udocs/cli/udocs"
)

// buildNumber is set via -ldflags
var buildNumber string

func main() {
	cmd.BuildNumber = buildNumber
	udocs.FetchAsset = Asset
	// commands MUST be in alphabetical order
	cmd.Root.AddCommand(
		cmd.Build(),
		cmd.Destroy(),
		cmd.Env(),
		cmd.Publish(),
		cmd.Pull(),
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
