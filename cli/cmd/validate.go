package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/UltimateSoftware/udocs/cli/udocs"
	"github.com/spf13/cobra"
)

func Validate() *cobra.Command {
	var validate = &cobra.Command{
		Use:   "validate",
		Short: "Validate a docs directory",
		Long:  `udocs-validate is a quick smoke-test that verifies the required contents of a docs directory.`,
		Run: func(cmd *cobra.Command, args []string) {
			abs, err := filepath.Abs(dir)
			if err != nil {
				fmt.Printf("Validation failed: unable to determine absolute path of docs directory: %v\n", err)
				os.Exit(-1)
			}

			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Validation failed: unable to determine current working directory: %v\n", err)
				os.Exit(-1)
			}

			if abs == cwd {
				fmt.Println("Validation failed: docs directory cannot be the current working directory")
				os.Exit(-1)
			}

			if err := udocs.Validate(dir); err != nil {
				fmt.Printf("Validation failed: %v\n", err)
				os.Exit(-1)
			}

			fmt.Println("Validation successful.")
		},
	}

	setFlag(validate, "dir")
	return validate
}
