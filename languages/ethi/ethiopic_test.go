package ethi

import "testing"

func TestVowelOrders(t *testing.T) {
	// Ha-series, all 7 vowel orders. 6th order ə is dropped at word
	// boundary (here: end of input), so ህ alone returns "h"; embed it
	// in a longer word to verify the schwa is preserved mid-word.
	cases := map[string]string{
		"ሀ":  "hä", "ሁ": "hu", "ሂ": "hi", "ሃ": "ha",
		"ሄ":  "he", "ሆ": "ho",
		"ህ":  "h",   // word-end ə dropped
		"ህለ": "həlä", // ə preserved mid-word
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAmharicSamples(t *testing.T) {
	// Common Amharic / biblical words. Our scheme uses ä for 1st-order
	// uniformly (including for alef/ayn) and drops schwa at word end.
	cases := map[string]string{
		"እግዚአብሔር": "ʾəgəziʾäbəḥer", // God (Igziabher; BGN/PCGN keeps medial ə)
		"ኢየሱስ":   "ʾiyäsus",       // Jesus
		"ሰላም":   "sälam",          // peace
		"አማርኛ":  "ʾämarəña",      // Amharic
		"ሰማይ":   "sämay",          // heaven
		"ምድር":   "mədər",          // earth
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDigits(t *testing.T) {
	if got := Transliterate("፩፪፫"); got != "123" {
		t.Errorf("digits: got %q, want %q", got, "123")
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{
		"Hello world",
		"\\v 1 In the beginning",
		"1:1",
	} {
		if got := Transliterate(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("እግዚአብሔር") {
		t.Error("expected true for Ethiopic")
	}
	if Contains("Hello world") {
		t.Error("expected false for Latin")
	}
}
