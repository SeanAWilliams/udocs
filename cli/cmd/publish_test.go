package cmd

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ultimatesoftware/udocs/cli/config"
	"github.com/ultimatesoftware/udocs/cli/server"
	"github.com/ultimatesoftware/udocs/cli/storage"
)

func TestPublish(t *testing.T) {
	settings := config.DefaultSettings()
	dao := storage.NewMockDao("")
	s := httptest.NewServer(server.New(&settings, dao))
	os.Setenv("UDOCS_PORT", s.Listener.Addr().String()[len("127.0.0.1:"):])
	defer s.Close()

	if err := runTestCommand(Publish(), "--dir ../../docs"); err != nil {
		t.Error(err)
	}

}
