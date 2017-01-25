package udocs

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/UltimateSoftware/udocs/cli/storage"
)

func TestValidate(t *testing.T) {
	dir, err := filepath.Abs("../../docs") // the actual docs directory for UDocs
	if err != nil {
		t.Fatal("udocs.TestValidate: UDocs docs directory is missing")
	}

	if err := Validate(dir); err != nil {
		t.Errorf("Validate(%s) => %v", dir, err)
	}
}

func TestBuild(t *testing.T) {
	dir, err := filepath.Abs("../../docs") // the actual docs directory for UDocs
	if err != nil {
		t.Fatal("udocs.TestBuild: UDocs docs directory is missing")
	}

	dao := storage.NewMockDao("/tmp")
	if err := Build("test-route", dir, dao); err != nil {
		t.Errorf("Build(test-route, %s, *APIMockDao) => %v", dir, err)
	}

	expectedFiles := []string{
		SIDEBAR_JSON,
		"/test-route/index.html",
		"/test-route/BestPractices.html",
	}

	for _, f := range expectedFiles {
		if _, err := dao.Fetch(f); err != nil {
			t.Errorf("Build(test-route, %s, *APIMockDao) => missing %s", dir, f)
		}
	}
}

const testSummary = `
# My Test 1.0 	(Route/Path)
* [Overview](README.md)
* [Tests](tests/README.md)
	* [Unit](test/unit.md)
`

func TestExtractRoute(t *testing.T) {
	expected := "my-test-1-0-route-path"
	summary := bytes.NewReader([]byte(testSummary))
	route := ExtractRoute(summary)
	if route != expected {
		t.Errorf("Expected: %s, Got: %s", expected, route)
	}
}
