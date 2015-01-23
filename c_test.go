package main

import (
	"reflect"
	"strings"
	"testing"
)

var (
	someC = `# 28 "/usr/include/stdint.h" 2 3 4
# 36 "/usr/include/stdint.h" 3 4
typedef signed char int8_t;
typedef short int int16_t;

enum
  {
          SCM_RIGHTS = 0x01





};
`
)

func TestTokenize(t *testing.T) {
	want := []string{
		"typedef",
		"signed",
		"char",
		"int8_t",
		";",
		"typedef",
		"short",
		"int",
		"int16_t",
		";",
		"enum",
		"{",
		"SCM_RIGHTS",
		"=",
		"0x01",
		"}",
		";",
	}
	have, err := cTokenize(strings.NewReader(someC))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("have: %#v, want: %#v", have, want)
	}
}
