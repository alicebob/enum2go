package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type cEnum struct {
	Name     string
	Elements []enumElem
}

type enumElem struct {
	Name  string
	Value []string
}

// ToEnum converts enum tokens to a cEnum object.
func ToEnum(toks []string) (*cEnum, error) {
	if toks[0] != "enum" {
		return nil, errors.New("not an enum token sequence")
	}
	toks = toks[1:]
	// Enums can be 'enum {...}' or 'enum enumname {...}'.
	name := ""
	if toks[0] != "{" {
		name = toks[0]
		toks = toks[1:]
	}

	if toks[0] != "{" {
		return nil, errors.New("expected an '{'")
	}
	if toks[len(toks)-1] != "}" {
		return nil, errors.New("expected to end in '}'")
	}
	toks = toks[1 : len(toks)-1]

	var elems []enumElem
	for _, es := range splitElem(toks) {
		// Cases we support:
		// enum {
		//   FOO = 42,
		//   BAR,
		//   BAZ = BAR,
		// }
		tokname := es[0]
		es = es[1:]
		if len(es) > 1 {
			if es[0] != "=" {
				return nil, errors.New("expected an '='")
			}
			es = es[1:]
		}
		elems = append(elems, enumElem{
			Name:  tokname,
			Value: es,
		})
	}
	return &cEnum{
		Name:     name,
		Elements: elems,
	}, nil
}

// splitElem splits the token sequence on all "," tokens.
func splitElem(toks []string) [][]string {
	var (
		s    [][]string
		elem []string
	)
	for _, t := range toks {
		if t == "," {
			if len(elem) > 0 {
				s = append(s, elem)
				elem = nil
			}
			continue
		}
		elem = append(elem, t)
	}
	if len(elem) > 0 {
		s = append(s, elem)
	}
	return s
}

// enumToGo knows Go.
func enumToGo(e *cEnum) (string, error) {
	var (
		b        = bytes.Buffer{}
		typeName = e.Name
	)

	b.WriteString(fmt.Sprintf("type %v int\n", typeName))
	b.WriteString("const (\n")

	i := 0
	values := map[string]int{} // C names -> value map
	for _, elem := range e.Elements {
		switch len(elem.Value) {
		case 0:
			// just the default, iota style.
		case 1:
			// Either simple constant or maybe a reference to another
			// enum value.
			t := elem.Value[0]
			// does it look like a number?
			if v, err := strconv.Atoi(t); err == nil {
				i = v
			} else {
				// is this a reference to a value we already had?
				if v, ok := values[t]; ok {
					i = v
				} else {
					return "", fmt.Errorf("unknown thing: %s", t)
				}
			}
		default:
			// nothing simple. This could be something such as '1<<12'. But
			// maybe it's not.
			return "", fmt.Errorf("value too complex: %s", strings.Join(elem.Value, " "))
		}
		values[elem.Name] = i
		eName := goName(elem.Name)
		b.WriteString(fmt.Sprintf("\t%v %v = %d\n", eName, typeName, i))
		i++
	}

	b.WriteString(")\n\n")

	return b.String(), nil
}

// TODO
func goName(n string) string {
	return n
}
