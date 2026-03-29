package obscure_test

import (
	"fmt"
	"io"
	"os"
	"strings"

	"vimagination.zapto.org/obscure"
)

func Example() {
	fmt.Print("Encoded: ")
	io.Copy(os.Stdout, obscure.NewEncoder([]byte("MY_KEY"), strings.NewReader("Hello, World"), false))

	fmt.Print("\nDecoded: ")
	io.Copy(os.Stdout, obscure.NewEncoder([]byte("MY_KEY"), strings.NewReader("Cjsse, Kersd"), true))

	// Output:
	// Encoded: Cjsse, Kersd
	// Decoded: Hello, World
}
