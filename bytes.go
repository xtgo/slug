// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slug

import (
	"encoding/hex"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// SlugBytes replaces each run of characters which are not unicode letters or
// numbers with a single hyphen, except for leading or trailing runs. Letters
// will be stripped of diacritical marks and lowercased. Letter or number
// codepoints that do not have combining marks or a lower-cased variant will
// be passed through unaltered.
func SlugBytes(s []byte) []byte {
	s = norm.NFKD.Bytes(s)
	buf := make([]byte, 0, len(s))
	dash := false
	for len(s) > 0 {
		r, i := utf8.DecodeRune(s)
		switch {
		// unicode 'letters' like mandarin characters pass through
		case unicode.IsOneOf(lat, r):
			buf = append(buf, s[:i]...)
			dash = true
		case unicode.IsOneOf(nop, r):
			// skip
		case dash:
			buf = append(buf, '-')
			dash = false
		}
		s = s[i:]
	}
	i := len(buf) - 1
	if i >= 0 && buf[i] == '-' {
		buf = buf[:i]
	}
	return buf
}

// SlugAsciiBytes is identical to SlugBytes, except that runs of one or more
// unicode letters or numbers that still fall outside the ASCII range will have
// their UTF-8 representation hex encoded and delimited by hyphens. As with
// SlugBytes, in no case will hyphens appear at either end of the returned
// string.
func SlugAsciiBytes(s []byte) []byte {
	s = norm.NFKD.Bytes(s)
	const m = utf8.UTFMax
	var (
		ib    [m * 3]byte
		ob    []byte
		buf   = make([]byte, 0, len(s))
		dash  = false
		latin = true
	)
	for len(s) > 0 {
		r, i := utf8.DecodeRune(s)
		switch {
		case unicode.IsOneOf(lat, r):
			r = unicode.ToLower(r)
			n := utf8.EncodeRune(ib[:m], r)
			if r >= 128 {
				if latin && dash {
					buf = append(buf, '-')
				}
				n = hex.Encode(ib[m:], ib[:n])
				ob = ib[m : m+n]
				latin = false
			} else {
				if !latin {
					buf = append(buf, '-')
				}
				ob = ib[:n]
				latin = true
			}
			dash = true
			buf = append(buf, ob...)
		case unicode.IsOneOf(nop, r):
			// skip
		case dash:
			buf = append(buf, '-')
			dash = false
			latin = true
		}
		s = s[i:]
	}
	i := len(buf) - 1
	if i >= 0 && buf[i] == '-' {
		buf = buf[:i]
	}
	return buf
}

// IsSlugAsciiBytes is equivalent to IsSlugAscii.
func IsSlugAsciiBytes(s []byte) bool {
	dash := true
	for _, b := range s {
		switch {
		case b == '-':
			if dash {
				return false
			}
			dash = true
		case 'a' <= b && b <= 'z', '0' <= b && b <= '9':
			dash = false
		default:
			return false
		}
	}
	return !dash
}
