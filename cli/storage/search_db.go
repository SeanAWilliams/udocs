package storage

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

type SearchDB struct {
	Path string
	mapping.IndexMapping
	bleve.Index
}

const (
	TITLE    = "title"
	MODIFIED = "modified"
)

func NewSearchDB(dir string) (*SearchDB, error) {
	os.RemoveAll(filepath.Dir(dir))
	os.MkdirAll(filepath.Dir(dir), 0755)
	indexMapping := buildIndexMapping()
	index, err := bleve.Open(dir)
	if err == bleve.ErrorIndexPathDoesNotExist {
		index, err = bleve.New(dir, indexMapping)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &SearchDB{
		Path:         dir,
		IndexMapping: indexMapping,
		Index:        index,
	}, nil
}

func buildIndexMapping() mapping.IndexMapping {
	textFieldAnalyzer := "en"
	pageMapping := bleve.NewDocumentMapping()

	enTextFieldMapping := bleve.NewTextFieldMapping()
	enTextFieldMapping.Analyzer = textFieldAnalyzer
	pageMapping.AddFieldMappingsAt(TITLE, enTextFieldMapping)

	dateTimeMapping := bleve.NewDateTimeFieldMapping()
	pageMapping.AddFieldMappingsAt(MODIFIED, dateTimeMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("page", pageMapping)
	indexMapping.DefaultAnalyzer = textFieldAnalyzer

	return indexMapping
}

type QueryResult struct {
	Phrase       string
	Total        uint64
	Took         float64
	QueryMatches []QueryMatch
}

func (qr *QueryResult) ToMap() map[string]interface{} {
	if qr == nil {
		return nil
	}

	return map[string]interface{}{
		"phrase":        qr.Phrase,
		"total":         qr.Total,
		"took":          qr.Took,
		"query_matches": qr.QueryMatches,
	}
}

type QueryMatch struct {
	ID       string
	Rank     int
	Score    float64
	Modified string
	Title    string
	Body     template.HTML
}

func (s *SearchDB) Query(phrase string) (*QueryResult, error) {
	sr := bleve.NewSearchRequest(bleve.NewQueryStringQuery(phrase))
	sr.Highlight = bleve.NewHighlightWithStyle("html")
	sr.Fields = []string{TITLE, MODIFIED}

	searchResults, err := s.Search(sr)
	if err != nil {
		return nil, err
	}

	qr := QueryResult{
		Phrase:       phrase,
		Took:         formatFloat(searchResults.Took.Seconds()),
		Total:        searchResults.Total,
		QueryMatches: make([]QueryMatch, 0),
	}

	for i, hit := range searchResults.Hits {
		qm := QueryMatch{
			Rank:     i + searchResults.Request.From + 1,
			ID:       hit.ID,
			Score:    hit.Score,
			Title:    hit.Fields[TITLE].(string),
			Modified: hit.Fields[MODIFIED].(string),
		}
		var buf bytes.Buffer
		for _, fragments := range hit.Fragments {
			for _, fragment := range fragments {
				buf.WriteString(fragment + "\n")
			}
		}
		qm.Body = template.HTML(buf.String())
		qr.QueryMatches = append(qr.QueryMatches, qm)
	}

	return &qr, nil
}

var htmlCharFilterRegexp = regexp.MustCompile(`</?[!\w]+((\s+\w+(\s*=\s*(?:".*?"|'.*?'|[^'">\s]+))?)+\s*|\s*)/?>`)

func filterHTMLTags(input []byte) []byte {
	return htmlCharFilterRegexp.ReplaceAllFunc(input, func(in []byte) []byte {
		return bytes.Repeat([]byte(``), len(in))
	})
}

func formatFloat(val float64) float64 {
	f, err := strconv.ParseFloat(fmt.Sprintf("%f", val), 64)
	if err != nil {
		return val
	}
	return f
}
