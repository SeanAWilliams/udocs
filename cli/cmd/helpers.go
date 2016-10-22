package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultimatesoftware/udocs/cli/udocs"
)

var (
	dir, git string
	headless bool
)

func setFlag(cmd *cobra.Command, flag string) {
	switch flag {
	case "dir":
		cmd.Flags().StringVarP(&dir, "dir", "d", "docs", "Directory containing your guide and markdown files")
	case "git":
		cmd.Flags().StringVarP(&git, "git", "g", "", "Git repository from which to pull docs from (i.e. <project>/<repo>)")
	case "headless":
		cmd.Flags().BoolVar(&headless, "headless", false, "Run UDocs server in headless mode")
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
