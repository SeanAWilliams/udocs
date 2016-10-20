package cmd

import "github.com/spf13/cobra"

var Root = &cobra.Command{
	Use:  "udocs",
	Long: `UDocs is a CLI library for Go that easily renders Markdown documentation guides to HTML, and serves them over HTTP.`,
}
