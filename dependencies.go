package main

import (
	"bytes"
	"regexp"
	"strings"
	"unsafe"
)

type DependencyScanner struct {
	text []byte
}

var re *regexp.Regexp

func init() {
	res := Map([]rune{'\'', '"', '`'}, func(q rune) string {
		s := string(q)
		return s + "[^" + s + "]*" + s
	})
	re = regexp.MustCompile(`(?m)\bnew\s+URL\s*\(\s*(` + strings.Join(res, "|") + `)\s*,\s*import\.meta\.url\s*(?:,\s*)?\)`)
}

func NewDependencyScanner(text []byte) *DependencyScanner {
	return &DependencyScanner{text}
}

func (s *DependencyScanner) String() string {
	return unsafe.String(&s.text[0], len(s.text))
}

func (s *DependencyScanner) Scan() []string {
	if !bytes.Contains(s.text, []byte("import.meta.url")) {
		return nil
	}
	matches := re.FindAllSubmatchIndex(s.text, -1)
	return Map(matches, func(match []int) string {
		beg, end := match[2]+1, match[3]-1 // +1/-1 to remove quotes
		return unsafe.String(&s.text[beg], end-beg)
	})
}
