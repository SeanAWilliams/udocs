package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ultimatesoftware/udocs/cli/config"
)

func Env() *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Show UDocs local environment information",
		Long:  `udocs-env lists the keys and values of all UDocs environment variables for the current user session.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.LoadSettings().String())
		},
	}
}
