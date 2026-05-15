package brahmic

import "testing"

// scriptCase pairs a Script with a small set of input→output expectations.
type scriptCase struct {
	sc    *Script
	cases map[string]string
}

func TestCoreSamples(t *testing.T) {
	all := []scriptCase{
		{Devanagari, map[string]string{
			"राम":          "rāma",
			"कृष्ण":         "kr̥ṣṇa",
			"नमस्ते":       "namaste",
			"भारत":         "bhārata",
			"१२३":           "123",
		}},
		{Bengali, map[string]string{
			"বাংলা":   "bāṁlā",
			"ঈশ্বর":   "īśbara", // strict ISO 15919 emits "ba" for ব
			"যীশু":    "yīśu",
		}},
		{Gurmukhi, map[string]string{
			"ਪੰਜਾਬੀ":  "paṁjābī",
			"ਸਤਿ":     "sati",
			"ਨਾਮ":     "nāma",
		}},
		{Gujarati, map[string]string{
			"ગુજરાતી":  "gujarātī",
			"ઈશ્વર":   "īśvara",
		}},
		{Oriya, map[string]string{
			"ଓଡ଼ିଆ":   "oṛiā",
			"ଈଶ୍ୱର":  "īśwara",
		}},
		{Tamil, map[string]string{
			"தமிழ்":   "tamiḻ",
			"இயேசு":   "iyēcu",
			"நாமம்":   "nāmam",
		}},
		{Telugu, map[string]string{
			"తెలుగు":  "telugu",
			"దేవుడు":  "dēvuḍu",
			"యేసు":    "yēsu",
		}},
		{Kannada, map[string]string{
			"ಕನ್ನಡ":   "kannaḍa",
			"ದೇವರು":  "dēvaru",
			"ಯೇಸು":    "yēsu",
		}},
		{Malayalam, map[string]string{
			"മലയാളം":  "malayāḷaṁ",
			"യേശു":    "yēśu",
		}},
		{Sinhala, map[string]string{
			"සිංහල":   "siṁhala",
			"දෙවියන්":  "deviyan",
		}},
		{Javanese, map[string]string{
			"ꦗꦮ":   "jawa",   // ja + wa
			"ꦏꦾ":   "kya",    // ka + pengkal (medial y) + inherent a
			"ꦏꦿ":   "kra",    // ka + cakra (medial r) + inherent a
			"ꦄꦏ꧀ꦱꦫ": "aksara", // a + ka + virama + sa + ra
		}},
		{Sundanese, map[string]string{
			"ᮞᮥᮔ᮪ᮓ": "sunda", // sa + u + na + virama + da
			"ᮘᮞ":   "basa",  // ba + sa
			"ᮊᮡ":   "kya",   // ka + pamingkal (medial y) + inherent a
		}},
		{Balinese, map[string]string{
			"ᬩᬮᬶ":   "bali",   // ba + la + i-sign
			"ᬓᬶᬢᬩ": "kitaba", // ka + i + ta + ba
		}},
		{Buginese, map[string]string{
			"ᨅᨔ":  "basa", // ba + sa
			"ᨈᨊ":  "tana", // ta + na
			"ᨕ":   "a",    // independent letter A
			"ᨕᨗ":  "i",    // letter A + sign I → independent /i/ via carrier
			"ᨕᨙ":  "e",    // letter A + sign E → independent /e/ via carrier
		}},
		{Batak, map[string]string{
			"ᯅᯖᯂ᯲": "batah", // ba + ta + ha + pangolat
			"ᯅᯖᯂ᯳": "batah", // same word with pangonangon (alternate virama)
			"ᯂᯂ":   "haha",  // ha + ha
		}},
	}
	for _, tc := range all {
		for in, want := range tc.cases {
			got := Transliterate(in, tc.sc)
			if got != want {
				t.Errorf("[%s] Transliterate(%q) = %q, want %q",
					tc.sc.Name, in, got, want)
			}
		}
	}
}

func TestDetect(t *testing.T) {
	cases := map[string]string{
		"ಆದಿಯಲ್ಲಿ":   "kannada",
		"राम":         "devanagari",
		"বাংলা":    "bengali",
		"தமிழ்":   "tamil",
		"തുടക്കം":  "malayalam",
		"":            "",
		"Hello world": "",
	}
	for in, want := range cases {
		sc := Detect(in)
		got := ""
		if sc != nil {
			got = sc.Name
		}
		if got != want {
			t.Errorf("Detect(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPassthrough(t *testing.T) {
	// Latin/ASCII content must round-trip unchanged regardless of script.
	for _, sc := range All {
		s := "\\v 1 Hello (world) 12:34"
		if got := Transliterate(s, sc); got != s {
			t.Errorf("[%s] passthrough mangled: %q", sc.Name, got)
		}
	}
}
