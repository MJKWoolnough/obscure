# obscure

[![CI](https://github.com/MJKWoolnough/obscure/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/byteio/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/obscure.svg)](https://pkg.go.dev/vimagination.zapto.org/byteio)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/obscure)](https://goreportcard.com/report/vimagination.zapto.org/byteio)

--
    import "vimagination.zapto.org/obscure"

Package obscure is used wrap a reader with an obscuring layer, using a cipher that is generated using a user-specified key.

Only letters and numbers are encoded, punctuation is not modified.

NB: This is not cryptographically secure and should not be used in scenarios where secrecy and/or privacy is important.

## Highlights

 - Simple interface to wrap a reader with a cipher that replaces letters and numbers.
 - Replaces all letters and number in the Unicode.L and Unicode.N ranges.
 - Only replaces characters with those that share the same UTF-8 byte size.

## Usage

```go
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"vimagination.zapto.org/obscure"
)

func main() {
	fmt.Print("Encoded: ")
	io.Copy(os.Stdout, obscure.NewEncoder([]byte("MY_KEY"), strings.NewReader("Hello, World"), false))

	fmt.Print("\nDecoded: ")
	io.Copy(os.Stdout, obscure.NewEncoder([]byte("MY_KEY"), strings.NewReader("Cjsse, Kersd"), true))

	// Output:
	// Encoded: Cjsse, Kersd
	// Decoded: Hello, World
}
```

## Command

A simple command is included in `cmd/obscure` that can be used to obscure text. The following are the flags:

```
  -d    decode file
  -i string
        input file (default "-")
  -k string
        key to generate seed from
  -o string
        output file (default "-")
```

Only the `-k` flag is necessary to obscure text.

## Git

Obscure can be used as a git filter by running the following commands:

```bash
git config filter.obscure.clean "obscure -k 'MY_KEY'"
git config filter.obscure.smudge "obscure -k 'MY_KEY' -d"
git config diff.unobscure.textconv cat
```

You can use the filename as part of the key by including the `%f` placeholder.

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/obscure
