Export enums from C to Go.


For when you need a bunch of enums from a C header file in Go.

## Status

It compiles. Ship it!

## Usage

You need to preprocess the header file (with `cpp`. Add `-I` flags for import
paths), and pipe that through `enum2go`. Currently only named enums are
supported. Give the names of the enums you want as arguments.

Example:

``` bash
cpp some/header.h | enum2go myenum anotherenum > gen_enums.go
```

Or in theory in .go code, but that doesn't work:
``` Go
//go:generate cpp some/header.h | enum2go -package $GOPACKAGE myenum anotherenum > gen_enums.go
```

## Caveats

Plenty of those.

 - The tokeniser is very simplistic, but it is good enough for the enums I need.
 - Enum values can be either empty (`iota` style), a simple integer, or refer to other value _in the same enum_. Anything else is currently not supported, although some would be easy to add (`1<<3` comes to mind)
 - Only named enums are supported
 - Bring Your Own `cpp`. This also breaks `go generate`, it seems.
