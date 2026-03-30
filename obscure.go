// Package obscure is used wrap a reader with an obscuring layer, using a cipher that is generated using a user-specified key.
package obscure

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"io"
	"math/rand"
	"slices"
	"unicode"
	"unicode/utf8"
)

// Encoder wraps an io.Reader and will replace letters and numbers using a
// simple cipher.
type Encoder struct {
	m    map[rune]rune
	r    *bufio.Reader
	skip int
}

// NewEncoder generates a new cipher using the given key, and wraps the Reader
// so it can be encoded.
//
// The decode flag, when set to true, will reverse the encoding.
func NewEncoder(key []byte, f io.Reader, decode bool) *Encoder {
	m := newCipherMap(key, decode)

	return &Encoder{m: m, r: bufio.NewReader(f)}
}

func newCipherMap(key []byte, flip bool) map[rune]rune {
	h := md5.Sum(key)
	r := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(h[:8]))))

	return cipher(r, flip, unicode.Lower, unicode.Upper, unicode.Number)
}

func cipher(r *rand.Rand, flip bool, sets ...*unicode.RangeTable) map[rune]rune {
	m := make(map[rune]rune)

	for _, set := range sets {
		cipherRanges16(m, r, flip, set.R16)
		cipherRanges32(m, r, flip, set.R32)
	}

	return m
}

func cipherRanges16(m map[rune]rune, r *rand.Rand, flip bool, ranges []unicode.Range16) {
	var chars [4][]rune

	for _, rng := range ranges {
		for c := rng.Lo; c < rng.Hi; c += rng.Stride {
			l := utf8.RuneLen(rune(c))

			chars[l-1] = append(chars[l-1], rune(c))
		}
	}

	shuffleToMap(m, r, flip, chars)
}

func cipherRanges32(m map[rune]rune, r *rand.Rand, flip bool, ranges []unicode.Range32) {
	var chars [4][]rune

	for _, rng := range ranges {
		for c := rng.Lo; c < rng.Hi; c += rng.Stride {
			l := utf8.RuneLen(rune(c))

			chars[l-1] = append(chars[l-1], rune(c))
		}
	}

	shuffleToMap(m, r, flip, chars)
}

func shuffleToMap(m map[rune]rune, r *rand.Rand, flip bool, splitChars [4][]rune) {
	for _, chars := range splitChars {
		for a, c := range shuffle(r, slices.Clone(chars)) {
			if flip {
				m[c] = chars[a]
			} else {
				m[chars[a]] = c
			}
		}
	}
}

func shuffle[T any](r *rand.Rand, s []T) []T {
	for i := len(s) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func (e *Encoder) Read(p []byte) (int, error) {
	var buf [4]byte

	q := p

	for len(q) > 0 {
		r, s, err := e.r.ReadRune()
		if err != nil {
			return len(p) - len(q), err
		}

		if c, ok := e.m[r]; ok {
			r = c
		}

		n := copy(q, buf[e.skip:utf8.EncodeRune(buf[:], r)])
		q = q[n:]
		e.skip = (e.skip + n) % s
	}

	if e.skip > 0 {
		e.r.UnreadRune()
	}

	return len(p), nil
}
