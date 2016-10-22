package cmd

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ultimatesoftware/udocs/cli/config"
	"github.com/ultimatesoftware/udocs/cli/server"
	"github.com/ultimatesoftware/udocs/cli/storage"
	"github.com/ultimatesoftware/udocs/cli/udocs"
)

func Serve() *cobra.Command {
	var serve = &cobra.Command{
		Use:   `serve`,
		Short: `Renders docs directories, and serves them locally over HTTP`,
		Long:  `udocs-serve renders given docs directories into static HTML files, and serves them over HTTP.`,
		Run: func(cmd *cobra.Command, args []string) {
			settings := config.LoadSettings()
			addr := settings.BindAddr + ":" + settings.Port

			var dao storage.Dao
			if url := settings.MongoURL; url != "" {
				dao = storage.NewMongoDBDao(url, udocs.SearchPath())
			} else {
				dao = storage.NewFileSystemDao(udocs.DeployPath(), 0755, udocs.SearchPath())
			}

			sidebar, _ := udocs.LoadSidebar(dao)

			if err := sidebar.Save(dao); err != nil {
				log.Fatalf("error: command.Serve: %v", err)
			}

			if _, ok := dao.(*storage.MongoDBDao); ok {
				for _, summary := range sidebar {
					if err := udocs.UpdateSearchIndex(summary, dao); err != nil {
						log.Printf("error: command.Serve: %v", err)
					}
				}
			}

			if headless {
				s := server.New(&settings, dao)
				fmt.Println(settings.String())
				log.Print("Running udocs-serve in headless mode")
				log.Printf("udocs is listening on %s:%s", settings.EntryPoint, settings.Port)
				log.Fatal(http.ListenAndServe(addr, s))
				return
			}

			settings.RootRoute = parseRoute()
			localServer := server.New(&settings, dao)

			if err := udocs.Build(settings.RootRoute, dir, dao); err != nil {
				log.Fatalf("error: command.Serve: %v\n", err)
			}

			go watchFiles(settings.RootRoute, dir, dao)
			abs, err := filepath.Abs(dir)
			if err != nil {
				abs = dir
			}
			log.Printf("Watching local directory %s for file changes", abs)

			log.Printf("Serving docs at http://localhost:%s/%s", settings.Port, settings.RootRoute)
			log.Printf("Press Ctrl-C to close when finished.")
			log.Fatal(http.ListenAndServe(addr, localServer))
		},
	}

	setFlag(serve, "dir")
	setFlag(serve, "headless")
	return serve
}

func watchFiles(route, dir string, dao storage.Dao) {
	watch, kill := make(chan struct{}, 0), make(chan error, 0)
	go udocs.WatchFiles(dir, watch, kill)
	for {
		select {
		case <-watch:
			if err := udocs.Build(route, dir, dao); err != nil {
				log.Fatalf("error: command.Serve: %v\n", err)
			}
		case err := <-kill:
			log.Fatalf("error: command.Serve: %v\n", err)
		}
	}
}
