//go:build ignore

// Extract the actual Markdown structure used by the website's renderer.
package main

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type structure struct {
	Headings []string `json:"headings"`
	Code     []string `json:"code"`
}

func main() {
	result := make(map[string]structure)
	for _, path := range os.Args[1:] {
		source, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		// Front matter is metadata, not part of the rendered document.
		if bytes.HasPrefix(source, []byte("---\n")) {
			if end := bytes.Index(source[4:], []byte("\n---")); end >= 0 {
				source = source[4+end+4:]
			}
		}
		// Marketing template calls render their YAML arguments into HTML before
		// Markdown parsing. Their indented quotations are not code examples.
		source = regexp.MustCompile("(?s)\\{\\{(?:pullquote|backgroundquote|quote|toolsblurbs|projects|books|libraries) `.*?`\\}\\}").ReplaceAll(source, nil)
		markdown := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID(), parser.WithHeadingAttribute()))
		document := markdown.Parser().Parse(text.NewReader(source))
		value := structure{Headings: []string{}, Code: []string{}}
		ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				return ast.WalkContinue, nil
			}
			switch node.Kind() {
			case ast.KindHeading:
				id, _ := node.AttributeString("id")
				value.Headings = append(value.Headings, string(id.([]byte)))
			case ast.KindCodeBlock, ast.KindFencedCodeBlock:
				var block bytes.Buffer
				for i := 0; i < node.Lines().Len(); i++ {
					line := node.Lines().At(i)
					block.Write(line.Value(source))
				}
				value.Code = append(value.Code, strings.TrimRight(block.String(), "\n"))
			}
			return ast.WalkContinue, nil
		})
		result[path] = value
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		panic(err)
	}
}
