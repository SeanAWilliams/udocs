package storage

import (
	"errors"
	"path/filepath"
	"strings"
)

// MockDao should only be used for testing purposes
type MockDao struct {
	Dao
	root  string
	pages map[string][]byte
}

func NewMockDao(root string) *MockDao {
	return &MockDao{root: root, pages: make(map[string][]byte)}
}

func (m *MockDao) Insert(id string, data []byte) error {
	m.pages[filepath.Join(m.root, id)] = data
	return nil
}

func (m *MockDao) Index(id, title string, data []byte) error {
	return nil
}

func (m *MockDao) Fetch(id string) ([]byte, error) {
	data, ok := m.pages[filepath.Join(m.root, id)]
	if !ok {
		return nil, errors.New(id + " not found")
	}
	return data, nil
}

func (m *MockDao) FetchGlob(pattern string) []string {
	var pages []string
	for id := range m.pages {
		if strings.HasPrefix(id, m.
			root+pattern) {
			pages = append(pages, id)
		}
	}
	return pages
}
