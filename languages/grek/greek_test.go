package grek

import "testing"

func TestBasicLetters(t *testing.T) {
	cases := map[string]string{
		"αβγδε":        "abgde",
		"ζηθικ":        "zēthik",
		"λμνξο":        "lmnxo",
		"πρστυ":        "prsty",
		"φχψω":         "phchpsō",
		"Λόγος":        "Logos",
		"Χριστός":      "Christos",
		"εἰρήνη":       "eirēnē",
		"final ς":      "final s",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRoughBreathing(t *testing.T) {
	cases := map[string]string{
		"Ἡρώδης":   "Hērōdēs",
		"Ἑβραῖος":  "Hebraios",
		"ὁ":        "ho",
		"ἡ":        "hē",
		"Ἰησοῦς":   "Iēsous", // smooth breathing — no h
		// υ/Υ after a/e/o/η forms a diphthong → "u".
		// Note: rough breathing on the 2nd vowel of a diphthong (e.g.
		// υἱός) places "h" inside the diphthong rather than before it;
		// this is a known limitation (see package doc).
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{
		"Hello world",
		"\\v 1 In the beginning",
	} {
		if got := Transliterate(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("Λόγος") {
		t.Error("expected true for Greek")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin")
	}
}
