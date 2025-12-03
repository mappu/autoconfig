package autoconfig

import (
	"testing"
)

func TestFormat(t *testing.T) {
	type testCase struct {
		input, expect string
	}

	cases := []testCase{
		// Simple case
		{"foo", "foo"},
		{"Foo", "Foo"},

		// Underscore style
		{"Foo_Bar", "Foo Bar"},
		{"__Foo_Bar__", "Foo Bar"},

		// Camelcase style
		{"FooBar", "Foo Bar"},
		{"FoBa", "Fo Ba"},
		{"FoB", "Fo B"},
		{"FBa", "F Ba"},
		{"FB", "FB"},
		{"TLSConfig", "TLS Config"},
		{"ConfigTLS", "Config TLS"},
		{"PrefixACRONYMSuffix", "Prefix ACRONYM Suffix"},
	}

	for _, tc := range cases {
		got := formatLabel(tc.input)
		if got != tc.expect {
			t.Errorf("formatLabel(%q): got %q, want %q", tc.input, got, tc.expect)
		}
	}
}
