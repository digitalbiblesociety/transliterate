package modi

import "testing"

func TestCoreSamples(t *testing.T) {
	cases := map[string]string{
		// 𑘎 KA → "ka" (inherent /a/)
		"𑘎": "ka",
		// 𑘎𑘿 KA + virama → "k"
		"𑘎𑘿": "k",
		// 𑘨𑘰𑘦 RA + ā-sign + MA → "rāma"
		"𑘨𑘰𑘦": "rāma",
		// 𑘎𑘽 KA + anusvara → "kaṁ"
		"𑘎𑘽": "kaṁ",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("𑘎") {
		t.Error("expected true for Modi string")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin string")
	}
}
