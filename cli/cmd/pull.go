package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/ultimatesoftware/udocs/cli/config"
)

func Pull() *cobra.Command {
	pull := &cobra.Command{
		Use:   "pull",
		Short: "Pull docs from remote Git repository",
		Long:  `udocs-pull pulls the docs directory from a remote Git repository.`,
		Run: func(cmd *cobra.Command, args []string) {
			if git == "" {
				fmt.Println("udocs-pull requires the --git flag to be set.\nRun `udocs destroy --help` for more information.")
				os.Exit(-1)
			}

			settings := config.LoadSettings()

			uri := fmt.Sprintf("%s:%s/api/%s", settings.EntryPoint, settings.Port, git)
			if err := pullDocs(uri); err != nil {
				fmt.Printf("Pull failed: %v\n", err)
				os.Exit(-1)
			}
			fmt.Printf("Successfully pulled guide for %s\n", git)
		},
	}

	setFlag(pull, "git")
	return pull
}

func pullDocs(uri string) error {
	resp, err := http.Post(uri, "text/plain", bytes.NewBuffer([]byte{}))
	if err != nil {
		return fmt.Errorf("udocs.Pull failed to POST to %s: %v", uri, err)
	}
	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("udocs.Pull was unable to read the HTTP response body: %v", err)
		}
		resp.Body.Close()
		return fmt.Errorf("udocs.Pull returned HTTP response: %s", string(body))
	}

	return nil
}
