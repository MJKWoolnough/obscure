package obscure

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type reader struct {
	io.Reader
}

type writer struct {
	io.Writer
}

func TestEncoder(t *testing.T) {
	for n, test := range [...]struct {
		Key, Input, Output string
	}{
		{"abc123", "Hello, World", "Vxffo, Lobfu"},
		{"abc1234", "Hello, World", "Shuur, Hraug"},
		{"ABC", "Lorem ipsum dolor sit amet consectetur adipiscing elit.\nDolor sit amet consectetur adipiscing elit quisque faucibus.", "Udpax efjnx ldwdp jeh sxah bdkjabhahnp slefejbekt aweh.\nMdwdp jeh sxah bdkjabhahnp slefejbekt aweh unejuna rsnbeqnj."},
		{"KEY", "ÀΩБλж ६๔፩５Ⅴ½⑩", "ÂȀԺŗӥ ⅼ⁰⑵꣐⑱١៵"},
	} {
		for m, bufSize := range [...]int{1, 2, 3, 100} {
			var encoded, decoded bytes.Buffer

			buf := make([]byte, bufSize)

			io.CopyBuffer(writer{&encoded}, NewEncoder([]byte(test.Key), reader{strings.NewReader(test.Input)}, false), buf)

			encodedStr := encoded.String()

			io.CopyBuffer(writer{&decoded}, NewEncoder([]byte(test.Key), reader{&encoded}, true), buf)

			decodedStr := decoded.String()

			if encodedStr != test.Output {
				t.Errorf("test %d.%d: expecting output %q, got %q", n+1, m+1, test.Output, encodedStr)
			} else if decodedStr != test.Input {
				t.Errorf("test %d.%d: expecting output %q, got %q", n+1, m+1, test.Input, decodedStr)
			}
		}
	}
}
