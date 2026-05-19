package mymr

import "testing"

func TestDefaultMLCTS2(t *testing.T) {
	cases := map[string]string{
		"မြန်မာ":     "mjanma",
		"ရန်ကုန်":    "jankoun",
		"မန္တလေး":    "manta.lei:",
		"မြို့":      "mjou.",
		"ဆရာ":       "hsa.ja",
		"ဆော":       "hso:",
		"ကြီး":       "kji:",
		"ဘုရား":     "bu.ja:",
		"ဗုဒ္ဓ":      "budda.",
		"စာအုပ်":    "saou",
		"အင်္ဂလိပ်": "inga.lei'",
		"၁၂၃၄၅":     "12345",
		"မြန်မာ၊ ရန်ကုန်။": "mjanma, jankoun.",
		"ဣ": "i.",
		"ဤ": "i",
		"ဥ": "u.",
		"ဦ": "u",
		"ဧ": "ei",
		"ဩ": "o:",
		"ဪ": "o",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIPA(t *testing.T) {
	cases := map[string]string{
		"မြန်မာ":     "mjàɴmà",
		"မြို့":      "mjo̰",
		"ရန်ကုန်":    "jàɴkòʊɴ",
		"ဆရာ":       "sʰa̰jà",
		"ကြီး":       "tɕí",
		"အင်္ဂလိပ်": "ʔɪ̀ɴɡa̰leɪʔ",
		"စာ":        "sà",
		"ဗုဒ္ဓ":      "bṵdda̰",
	}
	for in, want := range cases {
		if got := TransliterateIPA(in); got != want {
			t.Errorf("TransliterateIPA(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMLCTS(t *testing.T) {
	cases := map[string]string{
		"မြန်မာ":  "mranma",
		"မြို့":    "mrui.",
		"ဆရာ":    "hca.ra",
		"ဘုရား":  "bhu.ra:",
		"ဗုဒ္ဓ":   "buddha.",
		"ရန်ကုန်": "rankun",
		"ကြီး":    "kri:",
	}
	for in, want := range cases {
		if got := TransliterateMLCTS(in); got != want {
			t.Errorf("TransliterateMLCTS(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMedials(t *testing.T) {
	cases := []struct {
		in, want string
		mode     func(string) string
	}{
		{"ကျ", "kja.", Transliterate},
		{"ကျ", "tɕa̰", TransliterateIPA},
		{"ကျ", "kya.", TransliterateMLCTS},
		{"ကြ", "kja.", Transliterate},
		{"ကြ", "tɕa̰", TransliterateIPA},
		{"ကြ", "kra.", TransliterateMLCTS},
		{"ကွ", "kwa.", Transliterate},
		{"ကွ", "kwa̰", TransliterateIPA},
		{"လှ", "hla.", Transliterate},
		{"လှ", "l̥a̰", TransliterateIPA},
		{"လျှ", "hlja.", Transliterate},
		{"ဟျ", "sha.", Transliterate},
	}
	for _, c := range cases {
		if got := c.mode(c.in); got != c.want {
			t.Errorf("mode(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRhymeCoverage(t *testing.T) {
	cases := map[string]string{
		"ကာ":    "ka",
		"ကား":   "ka:",
		"ကာ့":   "ka.",
		"ကပ်":   "ka'",
		"ကိုက်": "kai'",
		"ကောက်": "kau'",
		"ကန်":   "kan",
		"ကန်း":  "kan:",
		"ကန့်":  "kan.",
		"ကေ":    "kei",
		"ကေး":   "kei:",
		"ကဲ":    "ke:",
		"ကင်":   "kin",
		"ကင်း":  "kin:",
		"ကင့်":  "kin.",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormalization(t *testing.T) {
	asatFirst := string([]rune{0x1000, 0x1014, 0x103A, 0x1037})
	aukFirst := string([]rune{0x1000, 0x1014, 0x1037, 0x103A})
	a := Transliterate(asatFirst)
	b := Transliterate(aukFirst)
	if a != b {
		t.Errorf("normalization: %q vs %q differ (%q vs %q)", asatFirst, aukFirst, a, b)
	}
}

func TestVirama(t *testing.T) {
	got := Transliterate("ဗုဒ္ဓ")
	if got == "bu.da.da." {
		t.Errorf("virama not closing first syllable: got %q", got)
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{"Hello", "1:1", "", "Burmese: မြန်မာ"} {
		got := Transliterate(s)
		if !Contains(s) && got != s {
			t.Errorf("non-Burmese %q: got %q, want passthrough", s, got)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("မြန်မာ") {
		t.Error("expected true for Burmese string")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin string")
	}
	if Contains("") {
		t.Error("expected false for empty string")
	}
}
