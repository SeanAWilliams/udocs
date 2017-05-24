package udocs

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shurcooL/github_flavored_markdown"
	"golang.org/x/net/html"
)

func processMarkdown(route string, data []byte) ([]byte, error) {
	dom, err := processDOM(filepath.Join("/", route), markdownToHTML(data))
	if err != nil {
		return nil, err
	}
	return stripDOM(dom, `<html><head></head><body>`, `</body></html>`), nil
}

func markdownToHTML(data []byte) []byte {
	return github_flavored_markdown.Markdown(data)
}

func processDOM(root string, htmlDoc []byte) ([]byte, error) {
	dom, err := html.Parse(bytes.NewReader(htmlDoc))
	if err != nil {
		return nil, err
	}

	var process func(*html.Node) ([]byte, error)
	process = func(node *html.Node) ([]byte, error) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "a":
				processAnchorElement(node, root)
			case "div":
				processDivElement(node)
			case "img":
				processImageElement(node, root)
			case "code":
				processCodeElement(node)
			case "table":
				processTableElement(node)
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if _, err := process(c); err != nil {
				return nil, err
			}
		}

		var buf bytes.Buffer
		if err := html.Render(&buf, node); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	return process(dom)
}

func processAnchorElement(node *html.Node, root string) {
	for i, a := range node.Attr {
		if a.Key == "href" && !strings.ContainsRune(a.Val, '#') && isMarkdownPage(a.Val) && !isRemoteURL(a.Val) {
			node.Attr[i].Val = getHTMLPath(root, a.Val)
		}
	}
}

func processImageElement(node *html.Node, root string) {
	for i, img := range node.Attr {
		if img.Key == "src" && !isRemoteURL(img.Val) {
			node.Attr[i].Val = getHTMLPath(root, img.Val)
		}
	}
}

func processCodeElement(node *html.Node) {
	if node.Attr == nil {
		node.Attr = []html.Attribute{
			html.Attribute{Key: "class", Val: "language-default"},
		}
	}
}

func processTableElement(node *html.Node) {
	node.Attr = append(node.Attr, html.Attribute{Key: "class", Val: "table"})
}

func processDivElement(node *html.Node) {
	for _, d := range node.Attr {
		if d.Key == "class" {
			// add syntax highlighting for code blocks (uses prism.js)
			if lang := strings.TrimPrefix(d.Val, "highlight highlight-"); lang != d.Val {
				langClass := fmt.Sprintf("language-%s", strings.ToLower(lang))
				if preElem := node.FirstChild; preElem != nil {
					codeElem := &html.Node{
						Type: html.ElementNode,
						Data: "code",
						Attr: []html.Attribute{html.Attribute{Key: "class", Val: langClass}},
					}
					if preElem.FirstChild != nil {
						codeElem.AppendChild(&html.Node{Type: html.TextNode, Data: preElem.FirstChild.Data})
					}
					node.FirstChild.FirstChild = codeElem
				}
			}
		}
	}
}

func stripDOM(dom []byte, prefix, suffix string) []byte {
	// strip outer HTML tags
	return dom[len(prefix) : len(dom)-len(suffix)]
}

func isMarkdownPage(pageID string) bool {
	return filepath.Ext(pageID) == ".md"
}

func isRemoteURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func getMarkdownPath(subpaths ...string) string {
	return convertPaths(INDEX_HTML, README_MD, ".html", ".md", subpaths...)
}

func getHTMLPath(subpaths ...string) string {
	return convertPaths(README_MD, INDEX_HTML, ".md", ".html", subpaths...)
}

func convertPaths(oldRoot, newRoot, oldExt, newExt string, subpaths ...string) string {
	if len(subpaths) == 0 {
		return ""
	}

	last := len(subpaths) - 1
	lastPath := subpaths[last]

	if base := filepath.Base(lastPath); strings.EqualFold(base, oldRoot) {
		lastPath = lastPath[:len(lastPath)-len(base)] + newRoot
	} else if filepath.Ext(lastPath) == oldExt {
		lastPath = lastPath[:len(lastPath)-len(oldExt)] + newExt
	}

	elems := make([]string, last)
	for i := range elems {
		elems[i] = subpaths[i]
	}
	elems = append(elems, lastPath)

	return filepath.Join(elems...)
}
