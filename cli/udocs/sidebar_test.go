package udocs

import (
	"testing"
)

var summary_md = []byte(`
# Test Summary

* [Overview](README.md)
* [Alpha](alpha/README.md)
	* [Sub-Alpha](alpha/sub-alpha.md)
		* [Sub-Sub-Alpha](alpha/sub-sub-alpha.md)
`)

func TestParseSummary(t *testing.T) {
	route := "test"
	expected := Summary{
		Route:  route,
		Header: "Test Summary",
		Pages: []Page{
			Page{
				Title:     "Overview",
				Path:      getHTMLPath(getPageID(route, "README.md")),
				TreeLevel: 1},
			Page{
				Title:     "Alpha",
				Path:      getHTMLPath(getPageID(route, "alpha/README.md")),
				TreeLevel: 1,
				SubPages: []Page{
					Page{
						Title:     "Sub-Alpha",
						Path:      getHTMLPath(getPageID(route, "alpha/sub-alpha.md")),
						TreeLevel: 2,
						SubPages: []Page{
							Page{
								Title:     "Sub-Sub-Alpha",
								Path:      getHTMLPath(getPageID(route, "alpha/sub-sub-alpha.md")),
								TreeLevel: 3}}}}}},
	}

	got, err := ParseSummary(route, summary_md)
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
		t.Log("Terminating test due to invalid invariant")
		return
	}

	if expected.Route != got.Route {
		t.Errorf("Route -> expected: %s got: %s", expected.Route, got.Route)
	}

	if expected.Header != got.Header {
		t.Errorf("Header -> expected: %s got: %s", expected.Header, got.Header)
	}

	if len(expected.Pages) != len(got.Pages) {
		t.Errorf("len(Pages) -> expected: %d got: %d", len(expected.Pages), len(got.Pages))
		t.Log("Terminating test due to failed test condition")
		return
	}

	for i := range expected.Pages {
		expectedPage, gotPage := expected.Pages[i], got.Pages[i]
		comparePages(t, expectedPage, gotPage)
	}

}

func comparePages(t *testing.T, expectedPage Page, gotPage Page) {
	if expectedPage.Title != gotPage.Title {
		t.Errorf("Page.Title -> expected: %s got: %s", expectedPage.Title, gotPage.Title)
	}

	if expectedPage.Path != gotPage.Path {
		t.Errorf("Page.Path -> expected: %s got: %s", expectedPage.Path, gotPage.Path)
	}

	if expectedPage.TreeLevel != gotPage.TreeLevel {
		t.Errorf("Page.TreeLevel -> expected: %d got: %d", expectedPage.TreeLevel, gotPage.TreeLevel)
	}

	for i := range expectedPage.SubPages {
		expectedSubPage, gotSubPage := expectedPage.SubPages[i], gotPage.SubPages[i]
		comparePages(t, expectedSubPage, gotSubPage)
	}
}
