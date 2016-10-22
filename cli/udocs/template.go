package udocs

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ultimatesoftware/udocs/static"
)

var (
	defaultTemplateDir   = "templates/v2"
	defaultTemplateFiles = []string{
		"document.html",
		"header.html",
		"navbar.html",
		"sidebar.html",
		"inner.html",
		"search.html",
	}
)

type Template struct {
	params map[string]interface{}
	files  []string
	tmpl   *template.Template
}

func LoadTemplateFiles(files ...string) ([]string, error) {
	tmp := filepath.Join(os.TempDir(), "udocs", defaultTemplateDir)
	os.MkdirAll(tmp, 0755)

	var tmpls []string
	for _, f := range files {
		data, err := static.Asset(filepath.Join(defaultTemplateDir, f))
		if err != nil {
			return []string{}, err
		}

		filename := filepath.Join(tmp, f)
		if err := ioutil.WriteFile(filename, data, 0755); err != nil {
			return []string{}, err
		}
		tmpls = append(tmpls, filename)
	}

	return tmpls, nil
}

func MustParseTemplate(params map[string]interface{}, files ...string) *Template {
	templates, err := LoadTemplateFiles(files...)
	if err != nil {
		panic(err)

	}

	tmpl := template.Must(template.ParseFiles(templates...))
	return &Template{
		params: params,
		files:  files,
		tmpl:   tmpl,
	}
}

func (t *Template) WithParameter(k string, v interface{}) *Template {
	if t.params == nil {
		t.params = make(map[string]interface{}, 1)
	}
	t.params[k] = v
	return t
}

func (t *Template) ExecuteTemplate(w io.Writer, name string, b []byte) error {
	html := template.HTML(string(b))
	data := struct {
		Content *template.HTML
		Params  map[string]interface{}
	}{
		Content: &html,
		Params:  t.params,
	}

	if err := t.tmpl.Lookup(name).Execute(w, data); err != nil {
		return fmt.Errorf("udocs.ExecuteTemplate: %v", err)
	}

	return nil
}

func DefaultTemplateFiles() []string {
	return defaultTemplateFiles
}
