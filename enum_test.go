package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	for src, want := range map[string]cEnum{
		"enum myattr { KEY1 = 0 , KEY2 = 2 , KEY3 = 3 , }": cEnum{
			Name: "myattr",
			Elements: []enumElem{
				{Name: "KEY1", Value: []string{"0"}},
				{Name: "KEY2", Value: []string{"2"}},
				{Name: "KEY3", Value: []string{"3"}},
			},
		},
		// no trailing ,
		"enum myattr2 { KEY1 = 0 , KEY2 = 2 , KEY3 = 3 }": cEnum{
			Name: "myattr2",
			Elements: []enumElem{
				{Name: "KEY1", Value: []string{"0"}},
				{Name: "KEY2", Value: []string{"2"}},
				{Name: "KEY3", Value: []string{"3"}},
			},
		},
		// missing values
		"enum myattr3 { KEY1 = 0 , KEY2 , KEY3 }": cEnum{
			Name: "myattr3",
			Elements: []enumElem{
				{Name: "KEY1", Value: []string{"0"}},
				{Name: "KEY2", Value: []string{}},
				{Name: "KEY3", Value: []string{}},
			},
		},
	} {
		tokens := strings.Fields(src)
		have, err := ToEnum(tokens)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(have, &want) {
			t.Errorf("have: %#v, want: %#v", have, &want)
		}
	}
}
