package udocs

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

var defaultTemplateFiles = []string{
	absTmplPath("v2/document.html"),
	absTmplPath("v2/header.html"),
	absTmplPath("v2/navbar.html"),
	absTmplPath("v2/sidebar.html"),
	absTmplPath("v2/inner.html"),
	absTmplPath("v2/search.html"),
}

type Template struct {
	params map[string]interface{}
	files  []string
	tmpl   *template.Template
}

func MustParseTemplate(params map[string]interface{}, files ...string) *Template {
	tmpl := template.Must(template.ParseFiles(files...))
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

func DefaultTemplateFiles(abs bool) []string {
	if abs {
		var tmpls []string
		for i, t := range defaultTemplateFiles {
			tmpls = append(tmpls, filepath.Join("static", t))
			fmt.Println(tmpls[i])
		}
		return tmpls
	}
	return defaultTemplateFiles
}

func absTmplPath(filename string) string {
	return filepath.Join(TemplatePath(), filename)
}
