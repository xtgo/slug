// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slug_test

import (
	"testing"

	"github.com/xtgo/slug"
)

func TestIsSlugAscii(t *testing.T) {
	tests := []struct {
		s string
		b bool
	}{
		{"", false},
		{"-", false},
		{"A", false},
		{"a", true},
		{"-a", false},
		{"a-", false},
		{"a-0", true},
		{"aa", true},
		{"a--0", false},
		{"abc世界def", false},
	}

	for _, test := range tests {
		if slug.IsSlugAscii(test.s) != test.b {
			t.Errorf("IsSlugAscii(%q) != %t", test.s, test.b)
		}
		if slug.IsSlugAsciiBytes([]byte(test.s)) != test.b {
			t.Errorf("IsSlugAsciiBytes(%q) != %t", test.s, test.b)
		}
	}
}

func TestSlugAscii(t *testing.T) {
	var tests = []struct{ in, out string }{
		{"", ""},
		{"ABC世界def-", "abc-e4b896e7958c-def"},
		{"012世界", "012-e4b896e7958c"},
		{"世界345", "e4b896e7958c-345"},
		{"012-世界-345", "012-e4b896e7958c-345"},
	}

	for _, test := range tests {
		if out := slug.SlugAscii(test.in); out != test.out {
			t.Errorf("SlugAscii: %q: %q != %q", test.in, out, test.out)
		}
		if out := slug.SlugAsciiBytes([]byte(test.in)); string(out) != test.out {
			t.Errorf("SlugAsciiBytes: %q: %q != %q", test.in, out, test.out)
		}
	}
}
