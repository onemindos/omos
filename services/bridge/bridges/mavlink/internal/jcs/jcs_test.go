package jcs

import "testing"

func TestCanonicalize(t *testing.T) {
	cases := []struct{ in, want string }{
		// Key order, whitespace, nesting.
		{`{ "b": 2, "a": {"d": [1, 2], "c": null} }`, `{"a":{"c":null,"d":[1,2]},"b":2}`},
		// Numbers in ES6 form (RFC 8785 Appendix B samples).
		{`[1.0, -0, 1e21, 1e-7, 0.000001, 333333333.33333329, 4.50, 2e-3, 1E30]`,
			`[1,0,1e+21,1e-7,0.000001,333333333.3333333,4.5,0.002,1e+30]`},
		// No HTML escaping; required escapes only; non-ASCII literal.
		{`{"s":"<a&b>é\n\u0001\"\\/"}`, "{\"s\":\"<a&b>é\\n\\u0001\\\"\\\\/\"}"},
		// UTF-16 ordering: U+FB33 (BMP, high) sorts after U+1F600 (surrogate pair 0xD83D).
		{"{\"דּ\":1,\"\U0001F600\":2}", "{\"\U0001F600\":2,\"דּ\":1}"},
	}
	for _, c := range cases {
		got, err := Canonicalize([]byte(c.in))
		if err != nil {
			t.Fatalf("%s: %v", c.in, err)
		}
		if string(got) != c.want {
			t.Errorf("Canonicalize(%s)\n got %s\nwant %s", c.in, got, c.want)
		}
	}
}

func TestCanonicalizeErrors(t *testing.T) {
	for _, in := range []string{`{`, `{"a":1} {}`, `[1e400]`} {
		if _, err := Canonicalize([]byte(in)); err == nil {
			t.Errorf("expected error for %s", in)
		}
	}
}
