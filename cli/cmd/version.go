package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var BuildNumber string

func Version() *cobra.Command {
	var version = &cobra.Command{
		Use:   "version",
		Short: "Show UDocs version",
		Long:  `udocs-version shows the version (build number) of the local UDocs installation.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("UDocs v%s\nCopyright, Ultimate Software 2016\nApache License, Version 2.0\n", BuildNumber)
		},
	}
	return version
}
