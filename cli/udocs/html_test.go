package udocs

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"golang.org/x/net/html"
)

const errFmt = "Given: %v\n\t  Expected: %v\n\t  Got: %v"

func TestMarkdownToHTML(t *testing.T) {
	given := "# Hello, world!\nThis is only a test."
	expected := `.*<h1>.*Hello, world!</h1>[[:space:]]*<p>This is only a test.</p>`
	got := markdownToHTML([]byte(given))

	match, err := regexp.Match(expected, got)
	if err != nil {
		t.Fatalf("Terminating test due to regexp failure: %v", err)
	}

	if !match {
		t.Errorf(errFmt, given, `.*<h1>.*Hello, world!</h1>[[:space:]]*<p>This is only a test.</p>`, string(got))
	}
}

func TestStripDOM(t *testing.T) {
	given := `<html><head></head><body><h1>Hello, world!</h1><p>This is only a test.</p></body></html>`
	prefix := `<html><head></head><body>`
	suffix := `</body></html>`
	expected := `<h1>Hello, world!</h1><p>This is only a test.</p>`
	got := stripDOM([]byte(given), prefix, suffix)

	if expected != string(got) {
		t.Errorf(errFmt, fmt.Sprintf("dom = %s\n\t\tprefix = %s\n\t\tsuffix = %s", given, prefix, suffix), expected, string(got))
	}
}

func extractNode(t *testing.T, node *html.Node) *html.Node {
	return node.FirstChild.FirstChild.NextSibling.FirstChild
}

func renderNode(t *testing.T, node *html.Node) string {
	var buf bytes.Buffer
	if err := html.Render(&buf, node); err != nil {
		t.Fatalf("Terminating test due to failed HTML render: %v", err)
	}

	dom := buf.Bytes()
	dom = dom[len(`<html><head></head><body>`) : len(dom)-len(`</body></html>`)]
	return string(dom)
}

func TestProcessAnchorElement(t *testing.T) {
	given, root := `<a href="anchor/page.md">Test</a>`, "/test"
	expected := `<a href="/test/anchor/page.html">Test</a>`

	node, err := html.Parse(bytes.NewReader([]byte(given)))
	if err != nil {
		t.Fatalf("Terminating test due to failed HTML parse: %v", err)
	}

	processAnchorElement(extractNode(t, node), root)
	if got := renderNode(t, node); expected != got {
		t.Errorf(errFmt, given, expected, got)
	}
}

func TestProcessImageElement(t *testing.T) {
	given, root := `<img src="image/pic.png"/>`, "/test"
	expected := `<img src="/test/image/pic.png"/>`

	node, err := html.Parse(bytes.NewReader([]byte(given)))
	if err != nil {
		t.Fatalf("Terminating test due to failed HTML parse: %v", err)
	}

	processImageElement(extractNode(t, node), root)
	if got := renderNode(t, node); expected != got {
		t.Errorf(errFmt, given, expected, got)
	}
}

func TestProcessImageElementWithRemoteSrc(t *testing.T) {
	given, root := `<img src="http://somesite.com/image/pic.png"/>`, "/test"
	expected := `<img src="http://somesite.com/image/pic.png"/>`

	node, err := html.Parse(bytes.NewReader([]byte(given)))
	if err != nil {
		t.Fatalf("Terminating test due to failed HTML parse: %v", err)
	}

	processImageElement(extractNode(t, node), root)
	if got := renderNode(t, node); expected != got {
		t.Errorf(errFmt, given, expected, got)
	}
}

func TestProcessTableElement(t *testing.T) {
	given := `<table></table>`
	expected := `<table class="table"></table>`

	node, err := html.Parse(bytes.NewReader([]byte(given)))
	if err != nil {
		t.Fatalf("Terminating test due to failed HTML parse: %v", err)
	}

	processTableElement(extractNode(t, node))
	if got := renderNode(t, node); expected != got {
		t.Errorf(errFmt, given, expected, got)
	}
}

func TestProcessCodeElement(t *testing.T) {
	given := `<code>{ if else then }</code>`
	expected := `<code class="language-default">{ if else then }</code>`

	node, err := html.Parse(bytes.NewReader([]byte(given)))
	if err != nil {
		t.Fatalf("Terminating test due to failed HTML parse: %v", err)
	}

	processCodeElement(extractNode(t, node))
	if got := renderNode(t, node); expected != got {
		t.Errorf(errFmt, given, expected, got)
	}
}

func TestProcessDivElement(t *testing.T) {
	given := `<div class="highlight highlight-default"><pre>{if else then}</pre><div>`
	expected := `<div class="highlight highlight-default"><pre><code class="language-default">{if else then}</code></pre><div></div></div>`

	node, err := html.Parse(bytes.NewReader([]byte(given)))
	if err != nil {
		t.Fatalf("Terminating test due to failed HTML parse: %v", err)
	}

	processDivElement(extractNode(t, node))
	if got := renderNode(t, node); expected != got {
		t.Errorf(errFmt, given, expected, got)
	}
}

func TestGetMarkdownPath(t *testing.T) {
	testCases := []struct {
		paths        []string
		markdownPath string
	}{
		{paths: []string{}, markdownPath: ""},
		{paths: []string{"/", "test", "page.md"}, markdownPath: "/test/page.md"},
		{paths: []string{"/", "test", "page.html"}, markdownPath: "/test/page.md"},
		{paths: []string{"/", "test", "image.png"}, markdownPath: "/test/image.png"},
		{paths: []string{"test", "page.md"}, markdownPath: "test/page.md"},
		{paths: []string{"test", "page.html"}, markdownPath: "test/page.md"},
		{paths: []string{"/", "test", "pages", "page.md"}, markdownPath: "/test/pages/page.md"},
		{paths: []string{"/", "test", "pages", "page.html"}, markdownPath: "/test/pages/page.md"},
		{paths: []string{"test", "index.html"}, markdownPath: "test/README.md"},
		{paths: []string{"/", "test", "index.html"}, markdownPath: "/test/README.md"},
	}

	for _, tc := range testCases {
		expected := tc.markdownPath
		got := getMarkdownPath(tc.paths...)
		if expected != got {
			t.Errorf(errFmt, tc.paths, expected, got)
		}
	}
}

func TestGetHTMLPath(t *testing.T) {
	testCases := []struct {
		paths    []string
		htmlPath string
	}{
		{paths: []string{}, htmlPath: ""},
		{paths: []string{"/", "test", "page.md"}, htmlPath: "/test/page.html"},
		{paths: []string{"/", "test", "page.html"}, htmlPath: "/test/page.html"},
		{paths: []string{"/", "test", "image.png"}, htmlPath: "/test/image.png"},
		{paths: []string{"test", "page.md"}, htmlPath: "test/page.html"},
		{paths: []string{"test", "page.html"}, htmlPath: "test/page.html"},
		{paths: []string{"/", "test", "pages", "page.md"}, htmlPath: "/test/pages/page.html"},
		{paths: []string{"/", "test", "pages", "page.html"}, htmlPath: "/test/pages/page.html"},
		{paths: []string{"test", "README.md"}, htmlPath: "test/index.html"},
		{paths: []string{"/", "test", "README.md"}, htmlPath: "/test/index.html"},
	}

	for _, tc := range testCases {
		expected := tc.htmlPath
		got := getHTMLPath(tc.paths...)
		if expected != got {
			t.Errorf(errFmt, tc.paths, expected, got)
		}
	}
}
