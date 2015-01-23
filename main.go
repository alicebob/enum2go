package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		header = flag.String("h", "", "header file (or STDIN)")
		packag = flag.String("package", "main", "package name")
	)
	flag.Parse()

	wantedEnums := map[string]struct{}{}
	for _, e := range flag.Args() {
		wantedEnums[e] = struct{}{}
	}

	var (
		fh = os.Stdin
	)

	if *header != "" {
		var err error
		if fh, err = os.Open(*header); err != nil {
			panic(err)
		}
	}

	tokens, err := cTokenize(fh)
	if err != nil {
		panic(err)
	}

	fmt.Printf("package %v\n\n// code generated by enum2go.\n\n", *packag)
	for _, toks := range enumTokens(tokens) {
		enum, err := ToEnum(toks)
		if err != nil {
			panic(err)
		}
		if _, ok := wantedEnums[enum.Name]; ok {
			src, err := enumToGo(enum)
			if err != nil {
				panic(err)
			}
			fmt.Print(src)
		}
	}
}
