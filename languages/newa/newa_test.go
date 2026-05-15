package newa

import "testing"

func TestCoreSamples(t *testing.T) {
	cases := map[string]string{
		// 𑐎 KA → "ka"
		"𑐎": "ka",
		// 𑐎𑑂 KA + virama → "k"
		"𑐎𑑂": "k",
		// 𑐬𑐵𑐩 RA + ā + MA → "rāma"
		"𑐬𑐵𑐩": "rāma",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("𑐎") {
		t.Error("expected true for Newa string")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin string")
	}
}
