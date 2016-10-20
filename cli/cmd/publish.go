package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
	"github.com/ultimatesoftware/udocs/cli/config"
	"github.com/ultimatesoftware/udocs/cli/udocs"
)

func Publish() *cobra.Command {
	publish := &cobra.Command{
		Use:   "publish",
		Short: "Publish docs to a remote UDocs host",
		Long:  `udocs-publish compresses and sends the docs directory to a remote UDocs server for hosting.`,
		Run: func(cmd *cobra.Command, args []string) {
			settings := config.LoadSettings()

			if err := udocs.Validate(dir); err != nil {
				fmt.Printf("Publish failed: %v\n", err)
				os.Exit(-1)
			}

			os.MkdirAll(filepath.Join(os.TempDir(), "udocs"), 0755)
			tarball := filepath.Join(os.TempDir(), "udocs", filepath.Base(dir)+".tar.gz")
			defer os.Remove(tarball)

			if err := archiver.TarGz(tarball, []string{dir}); err != nil {
				fmt.Printf("Publish failed: %v\n", err)
				os.Exit(-1)
			}

			tmp, err := os.Open(tarball)
			if err != nil {
				fmt.Printf("Publish failed: %v\n", err)
				os.Exit(-1)
			}

			route := parseRoute()
			uri := fmt.Sprintf("%s:%s/api/%s", settings.EntryPoint, settings.Port, route)
			if err := publishDocs(uri, tmp); err != nil {
				fmt.Printf("Publish failed: %v\n", err)
				os.Exit(-1)
			}
			fmt.Printf("Successfully published guide to %s\n", uri)
		},
	}

	setFlag(publish, "dir")
	return publish
}

// Publish sends an HTTP request to the server to publish the documentation in the build directory.
func publishDocs(uri string, r io.Reader) error {
	resp, err := http.Post(uri, "application/octet-stream", r)
	if err != nil {
		return fmt.Errorf("udocs.Publish failed to POST to %s: %v", uri, err)
	}

	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("udocs.Publish was unable to read the HTTP response body: %v", err)
		}
		resp.Body.Close()
		return fmt.Errorf("udocs.Publish returned HTTP response: %s", string(body))
	}

	return nil
}
