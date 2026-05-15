package tibt

import "testing"

func TestSamples(t *testing.T) {
	cases := map[string]string{
		"ཀ":     "ka",         // ka + inherent a
		"ཀོ":    "ko",         // ka + o
		"བོད":   "boda",       // bod: b+o, d+inherent a
		"ལྷ":   "lha",         // l + subjoined h + inherent a
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Logf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{"Hello", "\\v 1 In the beginning"} {
		if got := Transliterate(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}
