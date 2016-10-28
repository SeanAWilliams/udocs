package udocs

import (
	"html/template"
	"io"
	"log"

	rice "github.com/GeertJohan/go.rice"
)

type Template struct {
	params map[string]interface{}
	files  []string
	tmpl   *template.Template
}

func DefaultTemplateFiles() []string {
	return []string{
		"document.html",
		"header.html",
		"navbar.html",
		"sidebar.html",
		"inner.html",
		"search.html",
	}
}

func MustParseTemplate(params map[string]interface{}, files ...string) *Template {
	return &Template{
		params: params,
		files:  files,
		tmpl:   mustParseTemplate(params, files...),
	}
}

func mustParseTemplate(params map[string]interface{}, files ...string) *template.Template {
	box, err := rice.FindBox("../../static/templates/v2")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.New("")
	for _, f := range files {
		s, err := box.String(f)
		if err != nil {
			log.Fatal(err)
		}
		tmpl, err = tmpl.Parse(s)
		if err != nil {
			log.Fatal(err)
		}
	}

	return tmpl
}

func (t *Template) WithParameter(k string, v interface{}) *Template {
	if t.params == nil {
		t.params = make(map[string]interface{}, 1)
	}
	t.params[k] = v
	return t
}

func (t *Template) Execute(w io.Writer, name string, b []byte) error {
	html := template.HTML(string(b))
	data := struct {
		Content *template.HTML
		Params  map[string]interface{}
	}{
		Content: &html,
		Params:  t.params,
	}

	if err := t.tmpl.Lookup(name).Execute(w, data); err != nil {
		return err
	}

	return nil
}
