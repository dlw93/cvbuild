package main

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Document struct {
	root *html.Node
	refs map[atom.Atom]string
}

type VisitorFunc func(string) (string, error)

func NewDocument(r io.Reader) (*Document, error) {
	return newDocument(r, map[atom.Atom]string{
		atom.Script: "src",
		atom.Link:   "href",
		atom.Img:    "src",
		atom.A:      "href",
	})
}

func (d *Document) NewDocumentWithOptions(r io.Reader, targets map[string]string) (*Document, error) {
	if refs, err := lookup(targets); err != nil {
		return nil, err
	} else {
		return newDocument(r, refs)
	}
}

func (d *Document) Walk(v VisitorFunc) error {
	return d.walk(d.root, v)
}

func (d *Document) WriteTo(w io.Writer) error {
	return html.Render(w, d.root)
}

func newDocument(r io.Reader, refs map[atom.Atom]string) (*Document, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return &Document{root, refs}, nil
}

func (d *Document) walk(node *html.Node, v VisitorFunc) error {
	if node.Type == html.ElementNode {
		if name, ok := d.refs[node.DataAtom]; ok {
			for i, attr := range node.Attr {
				if attr.Key == name {
					if path, err := v(attr.Val); err != nil {
						return err
					} else {
						node.Attr[i] = html.Attribute{Key: name, Namespace: attr.Namespace, Val: path}
					}
					break
				}
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if err := d.walk(c, v); err != nil {
			return err
		}
	}

	return nil
}

func lookup(m map[string]string) (map[atom.Atom]string, error) {
	lcb := func(s string) []byte { return []byte(strings.ToLower(s)) }
	r := make(map[atom.Atom]string, len(m))
	for tag, attr := range m {
		if t := atom.Lookup(lcb(tag)); t == 0 {
			return nil, fmt.Errorf("invalid tag %q", tag)
		} else {
			r[t] = attr
		}
	}
	return r, nil
}
