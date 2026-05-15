package tibt

import "testing"

func run(t *testing.T, cases map[string]string) {
	t.Helper()
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBareConsonants(t *testing.T) {
	run(t, map[string]string{
		"ཀ": "ka",
		"ཁ": "kha",
		"ག": "ga",
		"ང": "nga",
		"ཨ": "a",
		"ས": "sa",
	})
}

func TestVowelSigns(t *testing.T) {
	// Achen ཨ is silent when carrying a vowel sign: ཨོ → "o" not "ao".
	run(t, map[string]string{
		"ཀི": "ki",
		"ཀུ": "ku",
		"ཀེ": "ke",
		"ཀོ": "ko",
		"ཀཻ": "kai",
		"ཀཽ": "kau",
		"ཨི": "i",
		"ཨུ": "u",
		"ཨེ": "e",
		"ཨོ": "o",
	})
}

func TestLongVowels(t *testing.T) {
	// F71+F72 / F71+F74 decomposed long vowels canonicalise to I/U.
	run(t, map[string]string{
		"ཀཱ":   "kA",
		"ཀཱི":  "kI",
		"ཀཱུ":  "kU",
		"ཀཱུཾ": "kUM",
	})
}

func TestSubjoinedStacks(t *testing.T) {
	run(t, map[string]string{
		"ལྷ":  "lha",
		"ཀྲ":  "kra",
		"བྲ":  "bra",  // b + subjoined r — contrast with བར below
		"རྒྱ": "rgya", // three-letter stack: r + subjoined g + subjoined y
	})
}

func TestPrefixRootSuffix(t *testing.T) {
	run(t, map[string]string{
		"བོད":  "bod",
		"བདེ":  "bde",
		"བཀ":   "bka",
		"བར":   "bar",
		"བཀར":  "bkar",
		"དཀར":  "dkar",
		"བཀྲ":  "bkra",
		"ཕྱག":  "phyag",
		"རྒྱལ":  "rgyal",
		"འཚལ":  "'tshal",
	})
}

func TestAchungParticle(t *testing.T) {
	// Achung + vowel is the genitive/instrumental particle, not a
	// root: པའི → "pa'i", not "p'i".
	run(t, map[string]string{
		"པའི":  "pa'i",
		"བའི":  "ba'i",
		"མཁའི": "mkha'i",
		"བོའི":  "bo'i",
		"སའི":  "sa'i",
		"པའུ":  "pa'u",
	})
}

func TestSecondSuffix(t *testing.T) {
	run(t, map[string]string{
		"ལེགས":    "legs",
		"མཁས":    "mkhas",
		"གྲགས":   "grags",
		"བསྒྲུབས": "bsgrubs",
	})
}

func TestDisambiguationDot(t *testing.T) {
	// "+" separates prefix+root pairs that would otherwise read as a
	// single subjoined cluster. EWTS uses "." but a period there would
	// confuse a TTS engine.
	run(t, map[string]string{
		"གཡོ":  "g+yo",
		"གྱོ":  "gyo",
		"གཡས":  "g+yas",
		"དཔའ":  "dpa'",
		"བརླ":  "b+rla",
		"བར":   "bar", // suffix r — explicit /a/ already disambiguates
		"དར":   "dar",
		"བདེ":  "bde",
	})
}

func TestSyllableSeparation(t *testing.T) {
	run(t, map[string]string{
		"བདེ་ལེགས":         "bde legs",
		"རྒྱལ་པོ":          "rgyal po",
		"ཤེས་རབ":          "shes rab",
		"ཐུབ་པ":           "thub pa",
		"བཀྲ་ཤིས":         "bkra shis",
		"བཀྲ་ཤིས་བདེ་ལེགས": "bkra shis bde legs",
	})
}

func TestMantra(t *testing.T) {
	// "padme" → "pdme" because the source has no tsek between p and
	// dme; strict Wylie puts /a/ only on the root stack (dm).
	run(t, map[string]string{
		"ཨོཾ་མ་ཎི་པདྨེ་ཧཱུྂ": "oM ma Ni pdme hUM",
	})
}

func TestPunctuation(t *testing.T) {
	run(t, map[string]string{
		"བདེ་ལེགས།": "bde legs,",
		"བདེ་ལེགས༎": "bde legs.",
		"ཀ།་ག།":     "ka, ga,",
	})
}

func TestCleanPunctuation(t *testing.T) {
	// "།།" (two adjacent shads, the common sentence-end pattern)
	// collapses to a period. No-break-tsek before shad doesn't leave
	// stray " ,". No doubled commas or periods anywhere.
	run(t, map[string]string{
		"བདེ་ལེགས།།": "bde legs.",
		"བདེ་ལེགས༎":  "bde legs.",
		"ཡིན༌ཞིང༌།":  "yin zhing,",
		"ཀ།། ག།།":    "ka. ga.",
	})
}

func TestDigits(t *testing.T) {
	run(t, map[string]string{
		"༡༢༣": "123",
		"༠༩":  "09",
	})
}

func TestAnusvaraVisarga(t *testing.T) {
	run(t, map[string]string{
		"ཀཾ":  "kaM",
		"ཀཿ":  "kaH",
		"ཨཾ":  "aM",
		"ཀོཾ": "koM",
	})
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{
		"Hello",
		"\\v 1 In the beginning",
		"123 abc",
	} {
		if got := Transliterate(s); got != s {
			t.Errorf("Transliterate(%q) = %q, want %q (passthrough)", s, got, s)
		}
	}
}

func TestEmpty(t *testing.T) {
	if got := Transliterate(""); got != "" {
		t.Errorf("Transliterate(\"\") = %q, want %q", got, "")
	}
}

func runPhonetic(t *testing.T, cases map[string]string) {
	t.Helper()
	for in, want := range cases {
		if got := TransliteratePhonetic(in); got != want {
			t.Errorf("TransliteratePhonetic(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPhoneticBareSyllables(t *testing.T) {
	runPhonetic(t, map[string]string{
		"ཀ":  "ka",
		"ག":  "ga",
		"པ":  "pa",
		"པོ": "po",
		"ཨ":  "a",
		"ཨོ": "o",
	})
}

func TestPhoneticSilentPrefix(t *testing.T) {
	runPhonetic(t, map[string]string{
		"བདེ":  "de",
		"དགེ":  "ge",
		"འདུན": "dün",
		"མདུན": "dün",
		"བཀྲ":  "tra",
	})
}

func TestPhoneticClusters(t *testing.T) {
	runPhonetic(t, map[string]string{
		"ཀྲ":  "tra",
		"ཁྲ":  "thra",
		"གྲ":  "dra",
		"པྱ":  "cha",
		"ཕྱ":  "cha",
		"བྱ":  "ja",
		"མྱ":  "nya",
		"ལྷ":  "lha",
		"ཟླ":  "da",
		"ཀླ":  "la",
		"རྒྱ": "gya",
		"སྒྱ": "gya",
	})
}

func TestPhoneticUmlauts(t *testing.T) {
	runPhonetic(t, map[string]string{
		"བོད":  "bö",
		"མན":   "men",
		"དགོན": "gön",
		"ལུས":  "lü",
		"སེམས": "sem",
	})
}

func TestPhoneticFinals(t *testing.T) {
	runPhonetic(t, map[string]string{
		"ཐུབ":  "thup",
		"ཕྱག":  "chak",
		"ལེགས": "lek",
		"ཤེར":  "sher",
		"དགོན": "gön",
	})
}

func TestPhoneticAggregateSamples(t *testing.T) {
	runPhonetic(t, map[string]string{
		"བདེ་ལེགས": "de lek",
		"བཀྲ་ཤིས":  "tra shi",
		"རྒྱལ་པོ":  "gyel po",
		"ཤེས་རབ":   "she rap",
		"ཕྱག་འཚལ":  "chak tshel",
		"དགེ་འདུན": "ge dün",
	})
}

func TestPhoneticPassthrough(t *testing.T) {
	for _, s := range []string{"Hello", "\\v 1 In the beginning"} {
		if got := TransliteratePhonetic(s); got != s {
			t.Errorf("TransliteratePhonetic(%q) = %q, want %q (passthrough)", s, got, s)
		}
	}
}

func TestSplitPhrases(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{
			"བདེ་ལེགས། རྒྱལ་པོ",
			[]string{"བདེ་ལེགས", "རྒྱལ་པོ"},
		},
		{
			"ཐོག་མར།  ས་གཞི།",
			[]string{"ཐོག་མར", "ས་གཞི"},
		},
		{
			"ཀ་ཁ་ག", // tsek isn't a phrase boundary
			[]string{"ཀ་ཁ་ག"},
		},
		{
			"",
			nil,
		},
		{
			"   ",
			nil,
		},
		{
			"hello, བོད world",
			[]string{"hello", "བོད", "world"},
		},
	}
	for _, c := range cases {
		got := SplitPhrases(c.in)
		if len(got) != len(c.want) {
			t.Errorf("SplitPhrases(%q) = %v (len %d), want %v (len %d)", c.in, got, len(got), c.want, len(c.want))
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("SplitPhrases(%q)[%d] = %q, want %q", c.in, i, got[i], c.want[i])
			}
		}
	}
}

func TestContains(t *testing.T) {
	cases := map[string]bool{
		"བོད":       true,
		"hello":     false,
		"hi བོད":    true,
		"":          false,
		"123":       false,
	}
	for in, want := range cases {
		if got := Contains(in); got != want {
			t.Errorf("Contains(%q) = %v, want %v", in, got, want)
		}
	}
}
