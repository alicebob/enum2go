package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	var (
		header      = flag.String("h", "", "header file to read (or STDIN)")
		outfile     = flag.String("o", "", "write Go code to file (or STDOUT)")
		packag      = flag.String("package", "main", "package name")
		cpp         = flag.String("cpp", "cpp", "C preprocessor command to use. Set to empty to disable.")
		cppIncludes = Strings{}
	)
	flag.Var(&cppIncludes, "I", "'-I' flags for cpp. Can be repeated.")
	flag.Parse()

	wantedEnums := map[string]struct{}{}
	for _, e := range flag.Args() {
		wantedEnums[e] = struct{}{}
	}

	var (
		fh io.ReadCloser = os.Stdin
	)

	if *header != "" {
		var err error
		if fh, err = os.Open(*header); err != nil {
			panic(err)
		}
	}

	if *cpp != "" {
		// run the header through cpp.
		args := []string{}
		for _, inc := range cppIncludes {
			args = append(args, "-I", inc)
		}
		cmd := exec.Command(*cpp, args...)
		cmd.Stdin = fh
		cmd.Stderr = os.Stderr

		var err error
		if fh, err = cmd.StdoutPipe(); err != nil {
			panic(err)
		}
		if err := cmd.Start(); err != nil {
			panic(err)
		}
	}

	tokens, err := cTokenize(fh)
	if err != nil {
		panic(err)
	}

	outfh := os.Stdout
	if *outfile != "" {
		var err error
		if outfh, err = os.Create(*outfile); err != nil {
			panic(err)
		}
		defer outfh.Close()
	}

	fmt.Fprintf(outfh, "package %v\n\n// code generated by enum2go.\n\n", *packag)
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
			fmt.Fprint(outfh, src)
		}
	}
}
