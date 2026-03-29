package obscure

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestEncoder(t *testing.T) {
	for n, test := range [...]struct {
		Key, Input, Output string
	}{
		{"abc123", "Hello, World", "Vxffo, Lobfu"},
		{"abc1234", "Hello, World", "Shuur, Hraug"},
		{"ABC", "Lorem ipsum dolor sit amet consectetur adipiscing elit.\nDolor sit amet consectetur adipiscing elit quisque faucibus.", "Udpax efjnx ldwdp jeh sxah bdkjabhahnp slefejbekt aweh.\nMdwdp jeh sxah bdkjabhahnp slefejbekt aweh unejuna rsnbeqnj."},
		{"KEY", "ÀΩБλж ६๔፩５Ⅴ½⑩", "ÂȀԺŗӥ ⅼ⁰⑵꣐⑱١៵"},
	} {
		var encoded, decoded bytes.Buffer

		io.Copy(&encoded, NewEncoder([]byte(test.Key), strings.NewReader(test.Input), false))

		encodedStr := encoded.String()

		io.Copy(&decoded, NewEncoder([]byte(test.Key), &encoded, true))

		decodedStr := decoded.String()

		if encodedStr != test.Output {
			t.Errorf("test %d: expecting output %q, got %q", n+1, test.Output, encodedStr)
		} else if decodedStr != test.Input {
			t.Errorf("test %d: expecting output %q, got %q", n+1, test.Input, decodedStr)
		}
	}
}
