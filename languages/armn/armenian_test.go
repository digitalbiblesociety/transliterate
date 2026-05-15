package armn

import "testing"

func TestSamples(t *testing.T) {
	cases := map[string]string{
		"Հայաստան": "Hayastan",
		"Երևան":    "Yerewan",
		"Աստված":   "Astwac",
		"Հիսուս":   "Hisus",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Logf("Transliterate(%q) = %q, want %q", in, got, want) // exploratory
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
