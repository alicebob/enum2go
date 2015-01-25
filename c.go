package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var (
	ops = []string{
		",", ";", "&", ",", "?", ":",
		"(", ")",
		"*=", "*",
		"|=", "|",
		"{", "}",
		"++", "+=", "+",
		"--", "-=", "-",
		"!=", "!",
		"==", "=",
		"<", "<<",
		">", ">>",
		"/=", "/",
		"%=", "%",
		"~=", "~",
		".", "^",
	}
)

// cTokenize tokenizes preprocessed C code.
func cTokenize(r io.Reader) ([]string, error) {
	var (
		b      = bufio.NewReader(r)
		tokens []string
	)
	var err error
	for err == nil {
		var l string
		l, err = b.ReadString('\n')
		if strings.HasPrefix(l, "#") {
			continue
		}
		tokens = append(tokens, tokenify(l)...)
	}
	if err == io.EOF {
		return tokens, nil
	}
	return nil, err
}

func tokenify(l string) []string {
	var ts []string
loop:
	for len(l) > 0 {

		// Operator?
		for _, op := range ops {
			if strings.HasPrefix(l, op) {
				ts = append(ts, op)
				l = l[len(op):]
				continue loop
			}
		}
		c := rune(l[0])

		// String?
		if c == '"' || c == '\'' {
			end := strings.IndexRune(l[1:], c)
			if end < 0 {
				end = len(l) - 1
			}
			ts = append(ts, l[:end+2])
			l = l[end+2:]
			continue
		}

		// Keyword or variable?
		isWord := func(r rune) bool {
			// TODO: this will do?
			return r >= '0' && r <= '9' ||
				r >= 'a' && r <= 'z' ||
				r >= 'A' && r <= 'Z' ||
				r == '_' ||
				r == '[' || // array lookups. yeah. well.
				r == ']'
		}
		if isWord(c) {
			end := strings.IndexFunc(l, func(r rune) bool { return !isWord(r) })
			if end < 0 {
				end = len(l)
			}
			ts = append(ts, l[:end])
			l = l[end:]
			continue
		}

		switch c {
		case ' ', '\r', '\n', '\t':
			// Whitespace?
			// Skip those.
			l = l[1:]
		default:
			// Tokeniser bug. No surprise there.
			fmt.Printf("unknown: %q\n", l[0])
			l = l[1:]
		}
	}
	return ts
}

// enumTokens finds all statements which start with 'enum' from the token
// stream.  We throw away everything else.
func enumTokens(toks []string) [][]string {
	var enums [][]string
	for len(toks) > 0 {
		if toks[0] != "enum" {
			toks = eatUpto(";", toks)
			continue
		}
		// enums always end with a }.
		for i, tok := range toks {
			if tok == "}" {
				enums = append(enums, toks[:i+1])
				toks = eatUpto(";", toks[i+1:])
				break
			}
		}
	}
	return enums
}

// eatUpto discards everything upto-and-including `upto`. It'll deal with
// nested {} blocks. The return is the rest of the token stream.
func eatUpto(upto string, toks []string) []string {
	for len(toks) > 0 {
		switch toks[0] {
		case upto:
			return toks[1:]
		case "{":
			toks = eatUpto("}", toks[1:])
			continue
		}
		toks = toks[1:]
	}
	return nil
}
