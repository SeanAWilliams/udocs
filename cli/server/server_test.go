package server

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ultimatesoftware/udocs/cli/config"
	"github.com/ultimatesoftware/udocs/cli/storage"
	"github.com/ultimatesoftware/udocs/cli/udocs"
)

func TestNew(t *testing.T) {
	settings := config.DefaultSettings()
	dao := storage.NewMockDao(os.TempDir())

	if s := New(&settings, dao); s == nil {
		t.Error("(s *server) cannot be nil")
	}
}

func TestHandle(t *testing.T) {
	settings := config.DefaultSettings()
	dao := storage.NewMockDao(udocs.DeployPath())
	server := New(&settings, dao)

	testData := []byte(`<h1>UDocs<\h1>`)
	dao.Insert("/udocs/index.html", testData)

	sidebar := make(udocs.Sidebar, 0)
	if err := sidebar.Save(dao); err != nil {
		log.Fatalf("error: command.Serve: %v\n", err)
	}

	w := bytes.NewBuffer([]byte{})
	tmpl := server.tmpl.WithParameter("sidebar", []udocs.Summary{udocs.Summary{}})
	if err := tmpl.ExecuteTemplate(w, "document", testData); err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	testServer := httptest.NewServer(server)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/udocs/index.html")
	if err != nil {
		t.Errorf("failed to execute GET: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET %s\tExpected: %d, Got: %d", resp.Request.URL, http.StatusOK, resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("failed to read response body")
	}

	if expected := w.String(); expected != string(data) {
		t.Errorf("GET %s\tExpected %s, Got: %s", resp.Request.URL, expected, string(data))
	}
}
