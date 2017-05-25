package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gdscheele/udocs/cli/config"
	"github.com/spf13/cobra"
)

func Destroy() *cobra.Command {
	destroy := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a docs directory from a remote UDocs server",
		Long: `
  udocs-destroy removes the contents of a docs directory from a remote UDocs server, but does not
  "un-register" that route from the server.
	`,
		Run: func(cmd *cobra.Command, args []string) {
			route := parseRouteFromSummary()
			settings := config.LoadSettings()
			uri := fmt.Sprintf("%s:%s/api/%s", settings.EntryPoint, settings.Port, route)

			if err := destroyDocs(uri); err != nil {
				fmt.Printf("Destroy failed: %v\n", err)
				os.Exit(-1)
			}

			fmt.Println("Successfully destroyed guide for " + route)
		},
	}

	setFlag(destroy, "dir")
	return destroy
}

func destroyDocs(uri string) error {
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("udocs.Destroy failed create HTTP request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("udocs.Destroy was failed to make the HTTP request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("udocs.Destroy was unable to read the HTTP response body: %v", err)
		}
		resp.Body.Close()
		return fmt.Errorf("udocs.Destroy returned HTTP response: %s", string(body))
	}

	return nil
}
