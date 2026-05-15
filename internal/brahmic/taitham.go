package brahmic

// Tai Tham (Lanna) (U+1A20..U+1AAF).
//
// Used for Northern Thai (Kam Mueang, ~6M speakers; Bible translations
// exist), Tai Lue / Xishuangbanna Dai, Khün, and as Lao Tham for
// Buddhist liturgical Pali in Laos. Codepoints derived from
// Aksharamukha (MIT), which maps the Tai Tham block to canonical
// Brahmic positions for ISO 15919-style romanization.
//
// Known simplifications:
//   - HIGH/LOW consonant class distinctions are flattened — tone class
//     matters for spoken Northern Thai but is not represented in plain
//     Latin output.
//   - Pre-positioned vowel taling (U+1A6E) is emitted in storage order
//     (before the consonant), not reordered. PyThaiNLP-style reordering
//     would need an engine extension.
//   - Compound vowel forms (e.g. AU = AA-letter + tarung + AA-sign) are
//     not detected as units — each component emits its own letter.
//   - Tone marks (U+1A75..U+1A7C) drop silently as in-block unmapped.
//   - Hora digits (U+1A80..U+1A89) and Tham digits (U+1A90..U+1A99)
//     both render as 0-9.
var TaiTham = func() *Script {
	sc := &Script{
		Name:             "lana",
		BlockStart:       0x1A20,
		BlockEnd:         0x1AAF,
		Virama:           []rune{0x1A7A}, // SAKOT
		DigitStart:       0x1A90,         // Tham digits
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}

	// Consonants in canonical Brahmic order using Aksharamukha's HIGH/LOW
	// selection. Voiceless aspirate slots (KH, CH, etc.) use HIGH class;
	// voiced aspirate slots (GH, JH, DH, BH) use LOW class.
	cons := map[rune]string{
		0x1A20: "k", 0x1A21: "kh", 0x1A23: "g", 0x1A25: "gh", 0x1A26: "ṅ",
		0x1A27: "c", 0x1A28: "ch", 0x1A29: "j", 0x1A2B: "jh", 0x1A2C: "ñ",
		0x1A2D: "ṭ", 0x1A2E: "ṭh", 0x1A2F: "ḍ", 0x1A30: "ḍh", 0x1A31: "ṇ",
		0x1A32: "t", 0x1A33: "th", 0x1A34: "d", 0x1A35: "dh", 0x1A36: "n",
		0x1A38: "p", 0x1A39: "ph", 0x1A3B: "b", 0x1A3D: "bh", 0x1A3E: "m",
		0x1A3F: "y", 0x1A41: "r",
		0x1A43: "l",
		0x1A45: "v",
		0x1A46: "ś", 0x1A47: "ṣ", 0x1A48: "s", 0x1A49: "h",
	}
	for r, v := range cons {
		sc.ConsonantBase[r] = v
	}

	// Independent vowels — single-rune subset. Compound forms (AA, R̥̄, AU)
	// are encoded multi-rune in Aksharamukha and emit component-by-component.
	sc.IndependentVowel[0x1A4B] = "a"
	sc.IndependentVowel[0x1A4D] = "i"
	sc.IndependentVowel[0x1A4E] = "ī"
	sc.IndependentVowel[0x1A4F] = "u"
	sc.IndependentVowel[0x1A50] = "ū"
	sc.IndependentVowel[0x1A51] = "e"
	sc.IndependentVowel[0x1A52] = "o"
	sc.IndependentVowel[0x1A42] = "r̥"
	sc.IndependentVowel[0x1A44] = "l̥"

	// Vowel signs (matras) — single-rune subset.
	sc.VowelSign[0x1A63] = "ā"  // tarung (long a)
	sc.VowelSign[0x1A65] = "i"  // mai kit
	sc.VowelSign[0x1A66] = "ī"  // mai sat
	sc.VowelSign[0x1A69] = "u"  // suku
	sc.VowelSign[0x1A6A] = "ū"  // suku riang
	sc.VowelSign[0x1A6E] = "e"  // taling (pre-positioned; not reordered)
	sc.VowelSign[0x1A6F] = "ae" // ae
	sc.VowelSign[0x1A70] = "o"  // oa (above)
	sc.VowelSign[0x1A71] = "ai" // ai
	sc.VowelSign[0x1A72] = "ay" // throng
	sc.VowelSign[0x1A73] = "ɨ"  // oa (below)

	sc.Special[0x1A74] = "ṁ" // mai kang (anusvara-like)

	// Hora digits — alternate digit set at U+1A80..U+1A89.
	for i := rune(0); i < 10; i++ {
		sc.Special[0x1A80+i] = string('0' + byte(i))
	}

	return sc
}()
