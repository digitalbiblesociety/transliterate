package brahmic

// Devanagari (U+0900..U+097F): Hindi, Marathi, Nepali, Sanskrit, etc.
//
// Notes:
//   - Devanagari has no short "e"/"o" at +0x0E/+0x12, but does have a
//     dedicated "short ā" (ऄ) at +0x04 in some texts; we ignore.
//   - The block adds nukta (़, +0x3C) which forms historically-Persian
//     consonants when combined: क़ ख़ ग़ ज़ ड़ ढ़ फ़ य़. These appear as a
//     two-rune sequence (consonant + nukta). We handle them by adding the
//     pre-composed codepoints (U+0958..U+095F, U+097B..U+097F) to the
//     consonant table; bare nukta is dropped.
//   - We do NOT implement Hindi-style schwa deletion. ISO 15919 preserves
//     inherent /a/ at all positions; that's the right default for Sanskrit
//     and Marathi. Hindi readers will see a slightly schwa-heavy form.
var Devanagari = func() *Script {
	sc := newStandardScript("devanagari", 0x0900)
	// Short e / short o do not exist in canonical Devanagari at +0x0E/+0x12.
	delete(sc.IndependentVowel, 0x090E)
	delete(sc.IndependentVowel, 0x0912)
	delete(sc.VowelSign, 0x0946)
	delete(sc.VowelSign, 0x094A)
	// Override +0x0F and +0x13 to plain e / o per ISO 15919 (Devanagari
	// has only one e and one o; "ē" / "ō" length doesn't apply).
	sc.IndependentVowel[0x090F] = "e"
	sc.IndependentVowel[0x0913] = "o"
	sc.VowelSign[0x0947] = "e"
	sc.VowelSign[0x094B] = "o"
	// Pre-composed nukta consonants.
	sc.ConsonantBase[0x0958] = "q"
	sc.ConsonantBase[0x0959] = "k͟h"
	sc.ConsonantBase[0x095A] = "ġ"
	sc.ConsonantBase[0x095B] = "z"
	sc.ConsonantBase[0x095C] = "ṛ"
	sc.ConsonantBase[0x095D] = "ṛh"
	sc.ConsonantBase[0x095E] = "f"
	sc.ConsonantBase[0x095F] = "ẏ"
	// ऍ ऎ ऑ — Devanagari diphthong vowels used for English loans.
	// ऍ (U+090D) → ê, ऑ (U+0911) → ô (informal but common).
	sc.IndependentVowel[0x090D] = "ê"
	sc.IndependentVowel[0x0911] = "ô"
	sc.VowelSign[0x0945] = "ê"
	sc.VowelSign[0x0949] = "ô"
	// Vedic OM (ॐ, U+0950)
	sc.Special[0x0950] = "oṁ"
	// Nukta combining (decomposed form: consonant + ़).
	sc.Nukta = 0x093C
	sc.NuktaCombine = map[rune]string{
		0x0915: "q",   // क़
		0x0916: "k͟h", // ख़
		0x0917: "ġ",   // ग़
		0x091C: "z",   // ज़
		0x0921: "ṛ",   // ड़
		0x0922: "ṛh",  // ढ़
		0x092B: "f",   // फ़
		0x092F: "ẏ",   // य़
	}
	return sc
}()
