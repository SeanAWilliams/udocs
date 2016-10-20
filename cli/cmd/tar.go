package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

func Tar() *cobra.Command {
	var tar = &cobra.Command{
		Use:   "tar",
		Short: "Tar a docs directory",
		Long:  `udocs-tar creates a docs.tar.gz file that can be sent to the server via an HTTP POST request.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := archiver.TarGz(filepath.Base(dir)+".tar.gz", []string{dir}); err != nil {
				fmt.Printf("Tar failed: %v\n", err)
				os.Exit(-1)
			}
		},
	}

	setFlag(tar, "dir")
	return tar
}
