package udocs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/UltimateSoftware/udocs/cli/storage"
	"github.com/fsnotify/fsnotify"
)

var innerHTMLTemplate = MustParseTemplate(nil, "inner.html")

const (
	README_MD    = "README.md"
	SUMMARY_MD   = "SUMMARY.md"
	SIDEBAR_JSON = "sidebar.json"
	INDEX_HTML   = "index.html"
)

// Validate validates the the given docs directory meets the format required by UDocs.
// Specifically, the directory must exist, be named "docs" , and include both a README.md and SUMMARY.md file at the root of the directory.
func Validate(dir string) error {
	if filepath.Base(dir) != "docs" {
		return fmt.Errorf("directory '%s' does not have a base path of 'docs'", dir)
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	rootLevelFiles, err := ioutil.ReadDir(abs)
	if err != nil {
		return err
	}

	if !containsFile(rootLevelFiles, README_MD) {
		return fmt.Errorf("missing %q file", README_MD)
	}

	if !containsFile(rootLevelFiles, SUMMARY_MD) {
		return fmt.Errorf("missing %q file", SUMMARY_MD)
	}

	return nil
}

func Build(route, dir string, dao storage.Dao) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	// cache all original pages, so we can mark any hanging pages for later deletion
	// purgeList := make(map[string]struct{}, 0)
	// for _, filename := range dao.FetchGlob(route + "/**/*") {
	// 	purgeList[filename] = struct{}{}
	// }

	var summary Summary
	foundSummary := false
	if err := filepath.Walk(abs, func(path string, fi os.FileInfo, err error) error {
		if fi == nil || !fi.Mode().IsRegular() {
			return nil
		}

		var id string

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if IsSummaryFile(path) && !foundSummary {
			foundSummary = true
			summary, err = ParseSummary(route, data)
			if err != nil {
				return err
			}

			sidebar, _ := LoadSidebar(dao) // load sidebar, if it exists

			data, err = json.Marshal(sidebar.Merge(summary))
			if err != nil {
				return err
			}
			id = SIDEBAR_JSON
		} else if isMarkdownPage(path) {
			data, err = processMarkdown(route, data)
			if err != nil {
				return err
			}

			buf := new(bytes.Buffer)
			if err := innerHTMLTemplate.Execute(buf, "inner", data); err != nil {
				return err
			}
			data = buf.Bytes()

			id = getHTMLPath("/", route, path[len(abs):])
		} else {
			id = filepath.Join("/", route, path[len(abs):])
		}

		if err := dao.Insert(id, data); err != nil {
			return err
		}

		// delete(purgeList, id)
		return nil

	}); err != nil {
		return err
	}

	// for filename := range purgeList {
	// 	if err := dao.Delete(filename); err != nil {
	// 		return err
	// 	}
	// }

	if err := LoadQuipDocuments(summary, dao); err != nil {
		return err
	}

	if err := UpdateSearchIndex(summary, dao); err != nil {
		return err
	}

	return nil
}

func LoadQuipDocuments(summary Summary, dao storage.Dao) error {
	var walk func(pages []Page) error
	walk = func(pages []Page) error {
		for _, page := range pages {
			id := filepath.Base(page.Path)
			if IsQuipThread(id) {
				thread, err := DefaultQuipClient.GetThread(strings.TrimSuffix(id, ".quip"))
				if err != nil {
					return err
				}

				if err := dao.Insert(getPageID(summary.Route, id), []byte(thread.HTML)); err != nil {
					return err
				}
			}

			if len(page.SubPages) > 0 {
				walk(page.SubPages)
			}
		}
		return nil
	}
	return walk(summary.Pages)
}

func WatchFiles(dir string, watch chan struct{}, kill chan error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		kill <- err
	}
	defer watcher.Close()

	if err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err := watcher.Add(path); err != nil {
			return err
		}
		return nil
	}); err != nil {
		kill <- err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod != fsnotify.Chmod {
				// only notify on changes to the actual file
				if strings.HasSuffix(event.Name, ".md") {
					log.Println("file modified: ", event.Name)
					watch <- struct{}{}
				}
			}
		case err := <-watcher.Errors:
			kill <- err
		}
	}
}

func ExtractRoute(r io.Reader) string {
	strs := strings.Fields(ParseSummaryHeader(bufio.NewScanner(r)))
	sb := new(bytes.Buffer)
	for i, s := range strs {
		sb.WriteString(strings.ToLower(s))
		if i != len(strs)-1 {
			sb.WriteRune('-')
		}
	}
	alphanum := regexp.MustCompile("[^a-z0-9]+")
	clean := regexp.MustCompile("^-+|-+$")
	return clean.ReplaceAllString(alphanum.ReplaceAllString(sb.String(), "-"), "")
}

func UpdateSearchIndex(summary Summary, dao storage.Dao) error {
	var walk func(pages []Page) error
	walk = func(pages []Page) error {
		for _, page := range pages {
			pageID, pageTitle := page.Path, page.Title
			if strings.HasSuffix(pageID, SIDEBAR_JSON) {
				continue
			}

			pageData, err := dao.Fetch(pageID)
			if err != nil {
				continue
			}

			if err := dao.Index(pageID, pageTitle, pageData); err != nil {
				return err
			}

			if len(page.SubPages) > 0 {
				walk(page.SubPages)
			}
		}
		return nil
	}
	return walk(summary.Pages)
}

func getPageID(route, path string) string {
	return filepath.Join("/", route, path)
}

func containsFile(files []os.FileInfo, name string) bool {
	for _, fi := range files {
		if fi.Mode().IsRegular() && fi.Name() == name {
			return true
		}
	}
	return false
}
