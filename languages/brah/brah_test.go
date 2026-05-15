package brah

import "testing"

func TestCoreSamples(t *testing.T) {
	cases := map[string]string{
		// 𑀭𑀸𑀫 (RA + ā-sign + MA) → "rāma"
		"𑀭𑀸𑀫": "rāma",
		// 𑀓 (KA) alone → "ka" via inherent /a/
		"𑀓": "ka",
		// 𑀓𑁆 (KA + virama) → "k" — virama suppresses inherent /a/
		"𑀓𑁆": "k",
		// 𑀓𑁆𑀱 (KA + virama + ṢA) → "kṣa" — conjunct with inherent /a/
		"𑀓𑁆𑀱": "kṣa",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSpecials(t *testing.T) {
	cases := map[string]string{
		// Anusvara
		"𑀓𑀁": "kaṁ",
		// Visarga
		"𑀓𑀂": "kaḥ",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDigits(t *testing.T) {
	// Modern decimal digits at U+11066..U+1106F
	if got := Transliterate("𑁦𑁧𑁨"); got != "012" {
		t.Errorf("digits: got %q", got)
	}
}

func TestContains(t *testing.T) {
	if !Contains("𑀓") {
		t.Error("expected true for Brahmi string")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin string")
	}
}

func TestPassthrough(t *testing.T) {
	if got := Transliterate("Hello"); got != "Hello" {
		t.Errorf("passthrough: got %q", got)
	}
}
