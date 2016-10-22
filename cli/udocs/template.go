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

func DefaultTemplateFiles() []string {
	return defaultTemplateFiles
}

func absTmplPath(filename string) string {
	return filepath.Join(TemplatePath(), filename)
}

var (
	DocumentTemplate = `
{{define "document"}}
<!DOCTYPE html>
<html lang="en">
  {{template "header" .}}
  <body>
    {{template "navbar" .}}
    <div class="container-fluid">
      <div id="parent" class="row">
      {{template "sidebar" .}}
      <div id="inner" class="col-sm-9 col-md-10 main">{{template "inner" .}}</div>
      </div>
    </div>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.2.0/jquery.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js"
        integrity="sha256-KXn5puMvxCw+dAYznun+drMdG1IFl3agK0p/pqT9KAo= sha512-2e8qq0ETcfWRI4HJBzQiA3UoyFk6tbNyG+qSaIBZLyW9Xf3sWZHN/lxe9fTh1U45DpPf07yj94KsUHHWe4Yk1A=="
        crossorigin="anonymous"></script>
    <script src='/static/scripts/app.js'></script>
    <script src="/static/scripts/prism.js"></script>
  </body>
</html>
{{end}}
	`

	HeaderTemplate = `
{{define "header"}}
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <meta name="HandheldFriendly" content="true"/>
    <meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-status-bar-style" content="black">

    <title>UDocs</title>

    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
    <link rel="stylesheet" href="/static/styles/prism.css">
    <!--<link rel="stylesheet" href="/static/styles/font-awesome.css">-->
    <link rel="stylesheet" href="/static/styles/app.css">

    <link href="https://fonts.googleapis.com/css?family=Ubuntu" rel="stylesheet">
    <link href="https://fonts.googleapis.com/css?family=Open+Sans:300,300i,400,400i" rel="stylesheet">
</head>
{{end}}
`

	InnerTemplate = `
{{define "inner"}} 
{{.Content}}
<br>
<hr/> {{if .Params.repo}}
<h4 style="text-align: center;">
	<a href='{{.Params.repo}}' rel="nofollow">
		<i class="fa fa-code"></i> View the source for this page in Stash <i class="fa fa-code"></i>
	</a>
</h4>
{{end}} 
{{end}}
`
	NavbarTemplate = `
{{define "navbar"}}
<nav class="navbar navbar-fixed-top">
      <div class="container-fluid">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          {{if ne .Params.organization ""}}
          <a class="navbar-brand" href="{{.Params.entrypoint}}">
                <span class="text">{{.Params.organization}} Documentation</span>
          </a>
          {{end}}
        </div>
        <div id="navbar" class="navbar-collapse collapse">
           <ul class="nav navbar-nav navbar-right">
           {{if ne .Params.email ""}}
            <li>
                <a href="mailto:{{.Params.email}}?Subject={{.Params.organization}}%20Docs%20Feedback">Feedback?</a>
            </li>
            {{end}}
            <li>
                <form id="navbar-search" class="navbar-form">
                    <div class="form-group">
                        <input id="search-input" type="search" class="form-control"
                               placeholder="{{.Params.search_placeholder}}">
                    </div>
                </form>
            </li>
        </ul>
        </div>
      </div>
    </nav>
{{end}}
`

	SearchTemplate = `
{{define "search"}}
<!DOCTYPE html>
<html>
{{template "header" .}}
<body>
{{template "navbar" .}}
<div class="container-fluid">
    <div id="parent" class="row">
    {{template "sidebar" .}}
    <div id="inner" class="col-sm-9 col-md-10 main">
        <div class="row"><h1>Search</h1></div>
            <div class="row"><h5>{{.Params.query_result.Total}} matches, took {{.Params.query_result.Took}} seconds</h5><hr></div>
            {{range .Params.query_result.QueryMatches}}
            <div class="row">
                <h4 style="margin-bottom: 0.25em;"><a title='{{.Title}}' href='{{.ID}}'>{{.Title}}</a></h4>
                <p style="font-size: 13px;"><code class="language-default">{{.ID}}</code><br>{{.Body}}</p>
            </div>
            {{end}}
        </div>
    </div>
</div>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/2.2.0/jquery.min.js"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js"
        integrity="sha256-KXn5puMvxCw+dAYznun+drMdG1IFl3agK0p/pqT9KAo= sha512-2e8qq0ETcfWRI4HJBzQiA3UoyFk6tbNyG+qSaIBZLyW9Xf3sWZHN/lxe9fTh1U45DpPf07yj94KsUHHWe4Yk1A=="
        crossorigin="anonymous"></script>
<script src='/static/scripts/app.js'></script>
</body>
</html>
{{end}}
`

	SidebarTemplate = `
{{define "sidebar"}}
<div id="main-sidebar-nav"  class="col-sm-3 col-md-2 sidebar">
    <ul class="nav-docs nav nav-sidebar">
    {{range .Params.sidebar}}{{if ne .Header ""}}
        <li class="main-item has-sub-items">
            <div class="has-sub-items-content" id="sidebar-main-text">{{.Header}}<i
                    class="fa fa-angle-down"></i></div>
            <ul class="sub-items">
                {{range .Pages}}
                {{if .SubPages}}
                <li class="has-sub-items">
                    <div class="has-sub-items-content"><a class="nav-docs-link" href='{{.Path}}'
                                                            title='{{.Title}}'>{{.Title}}</a><i
                            class="fa fa-angle-down"></i></div>
                    {{range .SubPages}}
                    {{if .SubPages}}
                    <ul class="sub-items level-2 has-sub-items">
                        <div class="has-sub-items-content"><a class="nav-docs-link" href='{{.Path}}'
                                                                title='{{.Title}}'>{{.Title}}</a><i
                                class="fa fa-angle-down"></i></div>
                        {{range .SubPages}}
                        {{if eq .TreeLevel 3}}
                        <ul class="sub-items level-3">
                            <li><a class="nav-docs-link" href='{{.Path}}' title='{{.Title}}'>{{.Title}}</a>
                            </li>
                        </ul>
                        {{else}}
                        <li><a class="nav-docs-link" href='{{.Path}}' title='{{.Title}}'>{{.Title}}</a></li>
                        {{end}}
                        {{end}}
                    </ul>
                    {{else}}
                    <ul class="sub-items level-2">
                        <li><a class="nav-docs-link" href='{{.Path}}' title='{{.Title}}'>{{.Title}}</a></li>
                    </ul>
                    {{end}}
                    {{end}}
                </li>
                {{else}}
                <li><a class="nav-docs-link" href='{{.Path}}' title='{{.Title}}'>{{.Title}}</a></li>
                {{end}}
                {{end}}
            </ul>
        </li>
    {{end}}{{end}}
    </ul>
</div>


<!--
<div id="main-sidebar-nav" class="sidebar">
    <nav class="nav-docs nav">
        <ul class="nav-docs-items">
            {{range .Params.sidebar}}
            {{if ne .Header ""}}
            <li class="main-item has-sub-items">
                <div class="has-sub-items-content" id="sidebar-main-text">{{.Header}}<i
                        class="fa fa-angle-down"></i></div>
                <ul class="sub-items">
                    {{range .Pages}}
                    {{if .SubPages}}
                    <li class="has-sub-items">
                        <div class="has-sub-items-content"><a class="nav-docs-link" href='{{.Path}}'
                                                              title='{{.Title}}'>{{.Title}}</a><i
                                class="fa fa-angle-down"></i></div>
                        {{range .SubPages}}
                        {{if .SubPages}}
                        <ul class="sub-items level-2 has-sub-items">
                            <div class="has-sub-items-content"><a class="nav-docs-link" href='{{.Path}}'
                                                                  title='{{.Title}}'>{{.Title}}</a><i
                                    class="fa fa-angle-down"></i></div>
                            {{range .SubPages}}
                            {{if eq .TreeLevel 3}}
                            <ul class="sub-items level-3">
                                <li><a class="nav-docs-link" href='{{.Path}}' title='{{.Title}}'>{{.Title}}</a>
                                </li>
                            </ul>
                            {{else}}
                            <li><a class="nav-docs-link" href='{{.Path}}' title='{{.Title}}'>{{.Title}}</a></li>
                            {{end}}
                            {{end}}
                        </ul>
                        {{else}}
                        <ul class="sub-items level-2">
                            <li><a class="nav-docs-link" href='{{.Path}}' title='{{.Title}}'>{{.Title}}</a></li>
                        </ul>
                        {{end}}
                        {{end}}
                    </li>
                    {{else}}
                    <li><a class="nav-docs-link" href='{{.Path}}' title='{{.Title}}'>{{.Title}}</a></li>
                    {{end}}
                    {{end}}
                </ul>
            </li>
            {{end}}
            {{end}}
        </ul>
    </nav>
</div>-->
{{end}}
`
)
