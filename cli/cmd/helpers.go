package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/UltimateSoftware/udocs/cli/udocs"
	"github.com/spf13/cobra"
)

var (
	dir             string
	headless, reset bool
)

func setFlag(cmd *cobra.Command, flag string) {
	switch flag {
	case "dir":
		cmd.Flags().StringVarP(&dir, "dir", "d", "docs", "Directory containing your guide and markdown files")
	case "headless":
		cmd.Flags().BoolVar(&headless, "headless", false, "Run UDocs server in headless mode")
	case "reset":
		cmd.Flags().BoolVar(&reset, "reset", false, "Reset local UDocs database")
	default:
		panic("command.setFlag: unrecognized flag --" + flag)
	}
}

func parseRoute() string {
	f, err := os.Open(filepath.Join(dir, udocs.SUMMARY_MD))
	if err != nil {
		return ""
	}
	defer f.Close()
	route := udocs.ExtractRoute(f)
	if route == "" {
		fmt.Printf("Failed to parse H1 header in SUMMARY.md, and no other route was specified.\nRun `udocs serve --help` for more information.")
		os.Exit(-1)
	}
	return route
}

func runTestCommand(cmd *cobra.Command, input string) error {
	cmd.SetArgs(strings.Split(input, " "))
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("failed to run test command `udocs %s %s` : %v", cmd.Name(), input, err)
	}
	return nil
}
