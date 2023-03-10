package main

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Document struct {
	root    *html.Node
	targets map[atom.Atom]string
}

type DependencyHandler func(string) (string, error)

func NewDocument(r io.Reader) (*Document, error) {
	return newDocument(r, map[atom.Atom]string{
		atom.Script: "src",
		atom.Link:   "href",
		atom.Img:    "src",
	})
}

func (d *Document) NewDocumentWithOptions(r io.Reader, targets map[string]string) (*Document, error) {
	if refs, err := lookup(targets); err != nil {
		return nil, err
	} else {
		return newDocument(r, refs)
	}
}

func (d *Document) Walk(h DependencyHandler) error {
	return d.walk(d.root, h)
}

func (d *Document) WriteTo(w io.Writer) error {
	return html.Render(w, d.root)
}

func newDocument(r io.Reader, targets map[atom.Atom]string) (*Document, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return &Document{root, targets}, nil
}

func (d *Document) walk(node *html.Node, h DependencyHandler) error {
	if node.Type == html.ElementNode {
		if name, ok := d.targets[node.DataAtom]; ok {
			for i, attr := range node.Attr {
				if attr.Key == name {
					if path, err := h(attr.Val); err != nil {
						return err
					} else {
						node.Attr[i].Val = path
					}
					break
				}
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if err := d.walk(c, h); err != nil {
			return err
		}
	}

	return nil
}

func lookup(m map[string]string) (map[atom.Atom]string, error) {
	r := make(map[atom.Atom]string, len(m))
	for tag, attr := range m {
		if t := atom.Lookup([]byte(tag)); t == 0 {
			return nil, fmt.Errorf("unknown tag %q", tag)
		} else {
			r[t] = attr
		}
	}
	return r, nil
}
