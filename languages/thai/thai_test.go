package thai

import "testing"

func TestPrepositionedReorder(t *testing.T) {
	// Pre-positioned vowels (เ แ โ ใ ไ) should appear AFTER their consonant.
	cases := map[string]string{
		"ไทย":   "thaiy",  // ไ + ท + ย — ai reordered to follow t-cluster
		"โลก":   "lok",    // โ + ล + ก
		"เมือง": "mueang", // เ + ม + ◌ื + อ + ง — uea vowel cluster
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMultiCharVowelPatterns(t *testing.T) {
	cases := map[string]string{
		"เกาะ": "ko",  // เCาะ → short o
		"เกา":  "kao", // เCา → ao diphthong
		"เกะ":  "ke",  // เCะ → short e
		"แปะ":  "pae", // แCะ → short ae
		"โต๊ะ":   "to", // โCะ with tone mark in the middle → short o
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSilentHo(t *testing.T) {
	// ห before a sonorant (ง น ม ย ร ล ว) is silent at syllable start.
	cases := map[string]string{
		"หมา":  "ma", // ห silent, ม + า
		"หนู":  "nu", // ห silent, น + ◌ู
		"หยา": "ya", // ห silent, ย + า
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestThanthakhat(t *testing.T) {
	// ◌์ silences the consonant beneath it. When the silenced consonant
	// is ร or ล (Sanskrit cluster pattern), the preceding consonant goes
	// silent too. Other silenced consonants are single-rune only.
	cases := map[string]string{
		"จันทร์":   "chan", // ทร์ collapses
		"เซ็นต์":  "sen",   // only ต silenced
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDigits(t *testing.T) {
	if got := Transliterate("๐๑๒๓๔๕๖๗๘๙"); got != "0123456789" {
		t.Errorf("digits: got %q", got)
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{"\\v 1 In the beginning", "Hello", "1:1"} {
		if got := Transliterate(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("ไทย") {
		t.Error("expected true for Thai string")
	}
	if Contains("Hello") {
		t.Error("expected false for non-Thai")
	}
}
