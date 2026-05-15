package shrd

import "testing"

func TestCoreSamples(t *testing.T) {
	cases := map[string]string{
		// 𑆑 KA → "ka"
		"𑆑": "ka",
		// 𑆑𑇀 KA + virama → "k"
		"𑆑𑇀": "k",
		// 𑆫𑆳𑆩 RA + ā + MA → "rāma"
		"𑆫𑆳𑆩": "rāma",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("𑆑") {
		t.Error("expected true for Sharada string")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin string")
	}
}
