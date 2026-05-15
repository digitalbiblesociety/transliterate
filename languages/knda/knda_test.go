package knda

import "testing"

func TestIndependentVowels(t *testing.T) {
	cases := map[string]string{
		"ಅ": "a", "ಆ": "ā", "ಇ": "i", "ಈ": "ī",
		"ಉ": "u", "ಊ": "ū", "ಎ": "e", "ಏ": "ē",
		"ಐ": "ai", "ಒ": "o", "ಓ": "ō", "ಔ": "au",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestInherentA(t *testing.T) {
	// Bare consonants emit consonant + inherent /a/.
	cases := map[string]string{
		"ಕ": "ka", "ಖ": "kha", "ಗ": "ga", "ಘ": "gha",
		"ಚ": "ca", "ಜ": "ja", "ಟ": "ṭa", "ಡ": "ḍa",
		"ತ": "ta", "ದ": "da", "ನ": "na", "ಪ": "pa",
		"ಬ": "ba", "ಮ": "ma", "ಯ": "ya", "ರ": "ra",
		"ಲ": "la", "ವ": "va", "ಶ": "śa", "ಷ": "ṣa",
		"ಸ": "sa", "ಹ": "ha", "ಳ": "ḷa",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestVowelSigns(t *testing.T) {
	// Consonant + vowel sign (matra) replaces the inherent /a/.
	cases := map[string]string{
		"ಕಾ": "kā", "ಕಿ": "ki", "ಕೀ": "kī", "ಕು": "ku",
		"ಕೂ": "kū", "ಕೆ": "ke", "ಕೇ": "kē", "ಕೈ": "kai",
		"ಕೊ": "ko", "ಕೋ": "kō", "ಕೌ": "kau", "ಕೃ": "kr̥",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestViramaClusters(t *testing.T) {
	// Virama (್) suppresses the inherent /a/, producing conjunct consonants.
	cases := map[string]string{
		"ಕ್": "k",       // bare consonant after virama
		"ಕ್ಷ": "kṣa",     // k + virama + ṣ + inherent a
		"ಕ್ಷಿ": "kṣi",    // k + virama + ṣ + i
		"ಸ್ಯ": "sya",     // s + virama + y + inherent a
		"ರಾಜ್ಯ": "rājya", // r+ā j+virama y+inherent a
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAnusvaraVisarga(t *testing.T) {
	cases := map[string]string{
		"ಅಂ":   "aṁ",
		"ಕಂ":   "kaṁ",
		"ಅಃ":   "aḥ",
		"ರಾಮಃ": "rāmaḥ",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNumerals(t *testing.T) {
	if got := Transliterate("೧೨೩೪೫೬೭೮೯೦"); got != "1234567890" {
		t.Errorf("got %q, want %q", got, "1234567890")
	}
}

func TestPassthrough(t *testing.T) {
	// Latin, punctuation, ASCII digits, whitespace should round-trip.
	for _, s := range []string{"Hello", "1:1", " ", "()[]{}", "verse 23"} {
		if got := Transliterate(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}

func TestBibleSamples(t *testing.T) {
	// Real lines from KANERV Genesis 1.
	cases := map[string]string{
		"ಆದಿಕಾಂಡ":       "ādikāṁḍa",          // Genesis (book title)
		"ಆದಿಯಲ್ಲಿ":      "ādiyalli",          // "in the beginning"
		"ದೇವರು":         "dēvaru",            // "God"
		"ಆಕಾಶ":          "ākāśa",             // "heaven/sky"
		"ಭೂಮಿ":          "bhūmi",             // "earth"
		"ಸೃಷ್ಟಿಸಿದನು":   "sr̥ṣṭisidanu",      // "created"
		"ಬೆಳಕು":         "beḷaku",            // "light"
		"ಭೂಲೋಕದ":        "bhūlōkada",         // "of the earth"
		"ಸಾಗರ":          "sāgara",            // "sea/ocean"
		"ಮೊದಲನೆ":        "modalane",          // "first"
		"ಚಲಿಸುತ್ತಿದ್ದನು": "calisuttiddanu", // "was moving/hovering"
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("ಆದಿಯಲ್ಲಿ") {
		t.Error("expected true for Kannada string")
	}
	if Contains("Hello world") {
		t.Error("expected false for non-Kannada string")
	}
	if Contains("verse 1: ಆದಿಯಲ್ಲಿ ದೇವರು") != true {
		t.Error("expected true for mixed string with Kannada")
	}
}
