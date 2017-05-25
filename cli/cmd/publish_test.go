package cmd

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gdscheele/udocs/cli/config"
	"github.com/gdscheele/udocs/cli/server"
	"github.com/gdscheele/udocs/cli/storage"
	"github.com/gdscheele/udocs/cli/udocs"
)

func TestPublish(t *testing.T) {
	settings := config.DefaultSettings()
	dao := storage.NewMockDao("")
	server.Tmpls = udocs.DefaultTemplateFiles()
	s := httptest.NewServer(server.New(&settings, dao))
	os.Setenv("UDOCS_PORT", s.Listener.Addr().String()[len("127.0.0.1:"):])
	defer s.Close()

	if err := runTestCommand(Publish(), "--dir ../../docs"); err != nil {
		t.Error(err)
	}

}
