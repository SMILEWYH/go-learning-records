//go:build ignore

// Compare documentation snippets with Go's lexer, retaining output assertions
// and compiler directives while allowing natural-language comments to change.
package main

import (
	"encoding/json"
	"go/scanner"
	"go/token"
	"os"
	"regexp"
	"strings"
)

type commentStructure struct {
	Tokens    []string `json:"tokens"`
	Protected []string `json:"protected"`
	Errors    []string `json:"errors"`
}

var outputMarker = regexp.MustCompile(`(?i)^//\s*(unordered\s+)?output:`)

func inspectComments(source string) commentStructure {
	result := commentStructure{Tokens: []string{}, Protected: []string{}, Errors: []string{}}
	file := token.NewFileSet().AddFile("snippet.go", -1, len(source))
	var lexer scanner.Scanner
	lexer.Init(file, []byte(source), func(_ token.Position, msg string) {
		result.Errors = append(result.Errors, msg)
	}, scanner.ScanComments)
	outputEnd := -1
	for {
		pos, tok, lit := lexer.Scan()
		if tok == token.EOF {
			break
		}
		if tok != token.COMMENT {
			result.Tokens = append(result.Tokens, tok.String()+":"+lit)
			// The scanner can emit an inserted semicolon between comments.
			if tok != token.SEMICOLON {
				outputEnd = -1
			}
			continue
		}
		offset := file.Offset(pos)
		continuation := outputEnd >= 0 && strings.TrimSpace(source[outputEnd:offset]) == ""
		if outputMarker.MatchString(lit) || continuation {
			result.Protected = append(result.Protected, lit)
			outputEnd = offset + len(lit)
		} else {
			outputEnd = -1
			if strings.HasPrefix(lit, "//go:") || strings.HasPrefix(lit, "//line ") || strings.HasPrefix(lit, "/*line ") {
				result.Protected = append(result.Protected, lit)
			}
		}
	}
	return result
}

func main() {
	var snippets map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&snippets); err != nil {
		panic(err)
	}
	result := make(map[string]commentStructure, len(snippets))
	for id, snippet := range snippets {
		result[id] = inspectComments(snippet)
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		panic(err)
	}
}
