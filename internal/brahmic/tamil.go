package brahmic

// Tamil (U+0B80..U+0BFF).
//
// Tamil is a significant outlier among Brahmic scripts:
//   - No script-level distinction between voiced/unvoiced or aspirated
//     consonants. The Tamil "k" stands for k, g, gh, kh depending on
//     phonological context; we render the script form only.
//   - "Grantha" letters (ஸ ஷ ஹ ஜ ஶ etc.) are imported for Sanskrit loans
//     and follow standard Brahmic offsets.
//   - SHA at U+0BB6 is present.
//   - No native digits widely used in modern Tamil text, but block has
//     them at +0xE6.
var Tamil = func() *Script {
	sc := &Script{
		Name:             "tamil",
		BlockStart:       0x0B80,
		BlockEnd:         0x0BFF,
		Virama:           []rune{0x0BCD},
		DigitStart:       0x0BE6,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}
	// Independent vowels.
	sc.IndependentVowel[0x0B85] = "a"
	sc.IndependentVowel[0x0B86] = "ā"
	sc.IndependentVowel[0x0B87] = "i"
	sc.IndependentVowel[0x0B88] = "ī"
	sc.IndependentVowel[0x0B89] = "u"
	sc.IndependentVowel[0x0B8A] = "ū"
	sc.IndependentVowel[0x0B8E] = "e"
	sc.IndependentVowel[0x0B8F] = "ē"
	sc.IndependentVowel[0x0B90] = "ai"
	sc.IndependentVowel[0x0B92] = "o"
	sc.IndependentVowel[0x0B93] = "ō"
	sc.IndependentVowel[0x0B94] = "au"
	// Consonants. Tamil has fewer than the standard 33; what it has:
	sc.ConsonantBase[0x0B95] = "k"  // க
	sc.ConsonantBase[0x0B99] = "ṅ"  // ங
	sc.ConsonantBase[0x0B9A] = "c"  // ச
	sc.ConsonantBase[0x0B9C] = "j"  // ஜ (grantha)
	sc.ConsonantBase[0x0B9E] = "ñ"  // ஞ
	sc.ConsonantBase[0x0B9F] = "ṭ"  // ட
	sc.ConsonantBase[0x0BA3] = "ṇ"  // ண
	sc.ConsonantBase[0x0BA4] = "t"  // த
	sc.ConsonantBase[0x0BA8] = "n"  // ந
	sc.ConsonantBase[0x0BA9] = "ṉ"  // ன
	sc.ConsonantBase[0x0BAA] = "p"  // ப
	sc.ConsonantBase[0x0BAE] = "m"  // ம
	sc.ConsonantBase[0x0BAF] = "y"  // ய
	sc.ConsonantBase[0x0BB0] = "r"  // ர
	sc.ConsonantBase[0x0BB1] = "ṟ"  // ற
	sc.ConsonantBase[0x0BB2] = "l"  // ல
	sc.ConsonantBase[0x0BB3] = "ḷ"  // ள
	sc.ConsonantBase[0x0BB4] = "ḻ"  // ழ
	sc.ConsonantBase[0x0BB5] = "v"  // வ
	sc.ConsonantBase[0x0BB6] = "ś"  // ஶ (grantha)
	sc.ConsonantBase[0x0BB7] = "ṣ"  // ஷ (grantha)
	sc.ConsonantBase[0x0BB8] = "s"  // ஸ (grantha)
	sc.ConsonantBase[0x0BB9] = "h"  // ஹ (grantha)
	// Vowel signs.
	sc.VowelSign[0x0BBE] = "ā"
	sc.VowelSign[0x0BBF] = "i"
	sc.VowelSign[0x0BC0] = "ī"
	sc.VowelSign[0x0BC1] = "u"
	sc.VowelSign[0x0BC2] = "ū"
	sc.VowelSign[0x0BC6] = "e"
	sc.VowelSign[0x0BC7] = "ē"
	sc.VowelSign[0x0BC8] = "ai"
	sc.VowelSign[0x0BCA] = "o"
	sc.VowelSign[0x0BCB] = "ō"
	sc.VowelSign[0x0BCC] = "au"
	// Special.
	sc.Special[0x0B82] = "ṁ" // anusvara (rarely used in Tamil)
	sc.Special[0x0B83] = "ḥ" // visarga / āytham
	return sc
}()
