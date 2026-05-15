package lana

import "testing"

func TestCoreSamples(t *testing.T) {
	cases := map[string]string{
		// ᨠ HIGH KA → "ka"
		"ᨠ": "ka",
		// ᨠ᩺ HIGH KA + SAKOT (virama) → "k"
		"ᨠ᩺": "k",
		// ᨠᩣᨾ KA + tarung (ā) + MA → "kāma"
		"ᨠᩣᨾ": "kāma",
		// ᩉᩴ MA + anusvara (mai kang) — anusvara directly after consonant
		"ᨠᩴ": "kaṁ",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestThamDigits(t *testing.T) {
	if got := Transliterate("᪀᪁᪂᪃᪄"); got != "01234" {
		t.Errorf("digits: got %q", got)
	}
}

func TestContains(t *testing.T) {
	if !Contains("ᨠ") {
		t.Error("expected true for Tai Tham string")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin string")
	}
}
