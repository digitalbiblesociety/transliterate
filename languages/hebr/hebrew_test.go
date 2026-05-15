package hebr

import "testing"

func TestUnpointedConsonants(t *testing.T) {
	cases := map[string]string{
		"שלום":   "šlwm",
		"ירושלים": "yrwšlym",
		"דוד":     "dwd",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPointedSamples(t *testing.T) {
	cases := map[string]string{
		"שָׁלוֹם":     "šālwōm",      // shalom (vav-holam, strict order)
		"בְּרֵאשִׁית": "bərēʾšiyt",  // bereshit (in the beginning)
		"אֱלֹהִים":   "ʾělōhiym",   // Elohim
		"דָּוִד":     "dāwid",
		"יִשְׂרָאֵל":  "yiśərāʾēl",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPunctuation(t *testing.T) {
	if got := Transliterate("שָׁלוֹם׃"); got != "šālwōm:" {
		t.Errorf("got %q", got)
	}
	if got := Transliterate("בֶּן־אָדָם"); got != "ben-ʾādām" {
		t.Errorf("got %q", got)
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{"Hello", "\\v 1 In the beginning"} {
		if got := Transliterate(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("שלום") {
		t.Error("expected true for Hebrew")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin")
	}
}
