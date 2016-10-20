package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var globalData *sync.RWMutex = new(sync.RWMutex)

type FileSystemDao struct {
	root string
	mode os.FileMode
	*SearchDB
}

func NewFileSystemDao(root string, mode os.FileMode, searchDir string) *FileSystemDao {
	return &FileSystemDao{
		root:     root,
		mode:     mode,
		SearchDB: NewSearchDB(searchDir),
	}
}

func (fs *FileSystemDao) Fetch(pageID string) ([]byte, error) {
	globalData.RLock()
	defer globalData.RUnlock()

	if filepath.Ext(pageID) == "" {
		pageID = filepath.Join(pageID, "index.html")
	}

	filename := filepath.Join(fs.root, pageID)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("storage.Fetch: %v", err)
	}
	return data, nil
}

func (fs *FileSystemDao) FetchGlob(pattern string) []string {
	globalData.RLock()
	defer globalData.RUnlock()

	ids := make([]string, 0)
	files, _ := filepath.Glob(filepath.Join(fs.root, pattern))
	for i, file := range files {
		if fi, err := os.Stat(file); err != nil || fi.IsDir() {
			continue
		}

		if filepath.Ext(file) != "" {
			files[i] = file[len(fs.root):]
			ids = append(ids, file[len(fs.root):])
		}

	}

	return ids
}

func (fs *FileSystemDao) Insert(pageID string, pageData []byte) error {
	globalData.Lock()
	defer globalData.Unlock()

	filename := filepath.Join(fs.root, pageID)
	if err := os.MkdirAll(filepath.Dir(filename), fs.mode); err != nil {
		return fmt.Errorf("storage.Insert: %v", err)
	}
	if err := ioutil.WriteFile(filename, pageData, fs.mode); err != nil {
		return fmt.Errorf("storage.Insert: %v", err)
	}

	return nil
}

func (fs *FileSystemDao) Delete(pageID string) error {
	globalData.Lock()
	defer globalData.Unlock()

	filename := filepath.Join(fs.root, pageID)
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("storage.Delete: %v", err)
	}

	if err := fs.SearchDB.Index.Delete(pageID); err != nil {
		return fmt.Errorf("storage.Delete: %v", err)
	}

	return nil
}

func (fs *FileSystemDao) DeleteGlob(pattern string) error {
	globalData.Lock()
	defer globalData.Unlock()

	files, err := filepath.Glob(filepath.Join(fs.root, pattern))
	if err != nil {
		return fmt.Errorf("storage.Delete: %v", err)
	}

	for _, f := range files {
		if err := os.RemoveAll(f); err != nil {
			log.Println(err.Error())
		}
		if err := fs.SearchDB.Index.Delete(f); err != nil {
			log.Println(err.Error())
		}
	}

	return nil
}

func (fs *FileSystemDao) Query(query string) (*QueryResult, error) {
	globalData.RLock()
	defer globalData.RUnlock()
	return fs.SearchDB.Query(query)
}

func (fs *FileSystemDao) Index(pageID, pageTitle string, pageData []byte) error {
	globalData.Lock()
	defer globalData.Unlock()

	indexData := struct {
		Title    string    `json:"title"`
		Body     string    `json:"body"`
		Modified time.Time `json:"modified"`
	}{
		Title:    pageTitle,
		Body:     string(filterHTMLTags(pageData)),
		Modified: time.Now(),
	}

	return fs.SearchDB.Index.Index(pageID, indexData)
}
