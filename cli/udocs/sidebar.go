package udocs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gdscheele/udocs/cli/storage"
)

type Sidebar []Summary

type Summary struct {
	Route  string `json:"route"`
	Header string `json:"header"`
	Pages  []Page `json:"pages"`
}

type Page struct {
	Title     string `json:"title"`
	Path      string `json:"path"`
	TreeLevel int    `json:"tree_level"`
	SubPages  []Page `json:"sub_pages"`
}

func LoadSidebar(dao storage.Dao) (Sidebar, error) {
	var sidebar Sidebar

	data, err := dao.Fetch(SIDEBAR_JSON)
	if err != nil {
		return sidebar, err
	}
	if err := json.Unmarshal(data, &sidebar); err != nil {
		return sidebar, err
	}

	return sidebar, nil
}

func (s Sidebar) Save(dao storage.Dao) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return dao.Insert(SIDEBAR_JSON, data)
}

func (s Sidebar) Merge(summary Summary) Sidebar {
	for i, item := range s {
		if item.Route == summary.Route {
			s[i] = summary
			return s
		}
	}

	if len(s) == 0 {
		s = make([]Summary, 0)
	}

	// when running udocs-serve locally, we may have arbitrary routes that are not predefined
	return append(s, summary)
}

func ParseSummaryHeader(scanner *bufio.Scanner) string {
	if scanner == nil {
		return ""
	}

	headerRegex := regexp.MustCompile(`# (.*)$`)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := headerRegex.FindStringSubmatch(line); matches != nil {
			return strings.TrimSpace(matches[len(matches)-1])
		}
	}

	return ""
}

func ParseSummary(route string, data []byte) (Summary, error) {
	summary := Summary{Route: route, Pages: make([]Page, 0)}
	scanner := bufio.NewScanner(bytes.NewReader(data))

	// parse the header
	if summary.Header = ParseSummaryHeader(scanner); summary.Header == "" {
		return summary, errors.New("udocs.ParseSummary did not find the H1 (header) line (i.e. '# My Guide')")
	}

	lineRegex := regexp.MustCompile(`\* \[(.*)\]\((.*)\)$`)

	// parse the remaining summary tree
	for scanner.Scan() {
		line := scanner.Text()
		indices := lineRegex.FindStringSubmatchIndex(line)
		// methods from package 'regexp' return slices with len() == 2 * (# of regex sub-expressions + 1),
		// and lineRegex has 2 sub-expressions, hence the magic number 6 below
		if len(indices) != 6 {
			// lineRegex did not match the 2 sub-expressions on this line, so we skip to next loop iteration
			continue
		}

		uri := line[indices[4]:indices[5]]
		page := Page{
			Title:     line[indices[2]:indices[3]],
			Path:      getHTMLPath(getPageID(route, uri)),
			TreeLevel: (indices[0] % 3) + 1,
		}

		i := len(summary.Pages) - 1
		switch page.TreeLevel {
		case 1:
			summary.Pages = append(summary.Pages, page)
		case 2:
			summary.Pages[i].SubPages = append(summary.Pages[i].SubPages, page)
		case 3:
			j := len(summary.Pages[i].SubPages) - 1
			if j < 0 {
				return summary, fmt.Errorf("udocs.ParseSummary failed to parse line: %s\n", line)
			}
			summary.Pages[i].SubPages[j].SubPages = append(summary.Pages[i].SubPages[j].SubPages, page)
		default:
			return summary, fmt.Errorf("udocs.ParseSummary failed to parse line %s (only 3 levels of pages are supported)\n", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return summary, fmt.Errorf("udocs.ParseSummary had a scanner error: %v\n", err)
	}

	return summary, nil
}

func IsSummaryFile(filename string) bool {
	return filepath.Base(filename) == SUMMARY_MD
}
