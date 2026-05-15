package tirh

import "testing"

func TestCoreSamples(t *testing.T) {
	cases := map[string]string{
		// 𑒏 KA → "ka"
		"𑒏": "ka",
		// 𑒏𑓂 KA + virama → "k"
		"𑒏𑓂": "k",
		// 𑒩𑒰𑒧 RA + ā + MA → "rāma"
		"𑒩𑒰𑒧": "rāma",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("𑒏") {
		t.Error("expected true for Tirhuta string")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin string")
	}
}
