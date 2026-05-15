package brahmic

// Sinhala (U+0D80..U+0DFF).
//
// Sinhala has a notably different block layout than the rest of the
// Brahmic family — the codepoints don't line up with the standard offsets
// at all. Written explicitly here.
//
// Sinhala features prenasalized consonants (ඟ ඬ ඳ ඹ) — single glyphs
// encoding a nasal followed by a homorganic stop. We render them as
// ⟨nasal + consonant⟩ in Latin using the combining breve-below (̆) to
// indicate prenasalization, per ISO 15919.
var Sinhala = func() *Script {
	sc := &Script{
		Name:             "sinhala",
		BlockStart:       0x0D80,
		BlockEnd:         0x0DFF,
		Virama:           []rune{0x0DCA}, // al-lakuna
		DigitStart:       0x0DE6,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}
	// Independent vowels.
	sc.IndependentVowel[0x0D85] = "a"   // අ
	sc.IndependentVowel[0x0D86] = "ā"   // ආ
	sc.IndependentVowel[0x0D87] = "ä"   // ඇ
	sc.IndependentVowel[0x0D88] = "ǟ"   // ඈ
	sc.IndependentVowel[0x0D89] = "i"   // ඉ
	sc.IndependentVowel[0x0D8A] = "ī"   // ඊ
	sc.IndependentVowel[0x0D8B] = "u"   // උ
	sc.IndependentVowel[0x0D8C] = "ū"   // ඌ
	sc.IndependentVowel[0x0D8D] = "r̥"  // ඍ
	sc.IndependentVowel[0x0D8E] = "r̥̄" // ඎ
	sc.IndependentVowel[0x0D8F] = "l̥"  // ඏ
	sc.IndependentVowel[0x0D90] = "l̥̄" // ඐ
	sc.IndependentVowel[0x0D91] = "e"   // එ
	sc.IndependentVowel[0x0D92] = "ē"   // ඒ
	sc.IndependentVowel[0x0D93] = "ai"  // ඓ
	sc.IndependentVowel[0x0D94] = "o"   // ඔ
	sc.IndependentVowel[0x0D95] = "ō"   // ඕ
	sc.IndependentVowel[0x0D96] = "au"  // ඖ
	// Consonants (codepoints per Unicode 15.1 Sinhala block).
	sc.ConsonantBase[0x0D9A] = "k"   // ක
	sc.ConsonantBase[0x0D9B] = "kh"  // ඛ
	sc.ConsonantBase[0x0D9C] = "g"   // ග
	sc.ConsonantBase[0x0D9D] = "gh"  // ඝ
	sc.ConsonantBase[0x0D9E] = "ṅ"   // ඞ
	sc.ConsonantBase[0x0D9F] = "n̆g" // ඟ prenasalized
	sc.ConsonantBase[0x0DA0] = "c"   // ච
	sc.ConsonantBase[0x0DA1] = "ch"  // ඡ
	sc.ConsonantBase[0x0DA2] = "j"   // ජ
	sc.ConsonantBase[0x0DA3] = "jh"  // ඣ
	sc.ConsonantBase[0x0DA4] = "ñ"   // ඤ
	sc.ConsonantBase[0x0DA5] = "ñj"  // ඥ jñya cluster
	sc.ConsonantBase[0x0DA6] = "n̆j" // ඦ prenasalized
	sc.ConsonantBase[0x0DA7] = "ṭ"   // ට
	sc.ConsonantBase[0x0DA8] = "ṭh"  // ඨ
	sc.ConsonantBase[0x0DA9] = "ḍ"   // ඩ
	sc.ConsonantBase[0x0DAA] = "ḍh"  // ඪ
	sc.ConsonantBase[0x0DAB] = "ṇ"   // ණ
	sc.ConsonantBase[0x0DAC] = "n̆ḍ" // ඬ prenasalized
	sc.ConsonantBase[0x0DAD] = "t"   // ත
	sc.ConsonantBase[0x0DAE] = "th"  // ථ
	sc.ConsonantBase[0x0DAF] = "d"   // ද
	sc.ConsonantBase[0x0DB0] = "dh"  // ධ
	sc.ConsonantBase[0x0DB1] = "n"   // න
	sc.ConsonantBase[0x0DB3] = "n̆d" // ඳ prenasalized
	sc.ConsonantBase[0x0DB4] = "p"   // ප
	sc.ConsonantBase[0x0DB5] = "ph"  // ඵ
	sc.ConsonantBase[0x0DB6] = "b"   // බ
	sc.ConsonantBase[0x0DB7] = "bh"  // භ
	sc.ConsonantBase[0x0DB8] = "m"   // ම
	sc.ConsonantBase[0x0DB9] = "m̆b" // ඹ prenasalized
	sc.ConsonantBase[0x0DBA] = "y"   // ය
	sc.ConsonantBase[0x0DBB] = "r"   // ර
	sc.ConsonantBase[0x0DBD] = "l"   // ල
	sc.ConsonantBase[0x0DC0] = "v"   // ව
	sc.ConsonantBase[0x0DC1] = "ś"   // ශ
	sc.ConsonantBase[0x0DC2] = "ṣ"   // ෂ
	sc.ConsonantBase[0x0DC3] = "s"   // ස
	sc.ConsonantBase[0x0DC4] = "h"   // හ
	sc.ConsonantBase[0x0DC5] = "ḷ"   // ළ
	sc.ConsonantBase[0x0DC6] = "f"   // ෆ
	// Vowel signs (matras).
	sc.VowelSign[0x0DCF] = "ā"   // ා
	sc.VowelSign[0x0DD0] = "ä"   // ැ
	sc.VowelSign[0x0DD1] = "ǟ"   // ෑ
	sc.VowelSign[0x0DD2] = "i"   // ි
	sc.VowelSign[0x0DD3] = "ī"   // ී
	sc.VowelSign[0x0DD4] = "u"   // ු
	sc.VowelSign[0x0DD6] = "ū"   // ූ
	sc.VowelSign[0x0DD8] = "r̥"  // ෘ
	sc.VowelSign[0x0DDF] = "l̥"  // ෟ
	sc.VowelSign[0x0DF2] = "r̥̄" // ෲ
	sc.VowelSign[0x0DF3] = "l̥̄" // ෳ
	sc.VowelSign[0x0DD9] = "e"   // ෙ
	sc.VowelSign[0x0DDA] = "ē"   // ේ
	sc.VowelSign[0x0DDB] = "ai"  // ෛ
	sc.VowelSign[0x0DDC] = "o"   // ො
	sc.VowelSign[0x0DDD] = "ō"   // ෝ
	sc.VowelSign[0x0DDE] = "au"  // ෞ
	// Special.
	sc.Special[0x0D82] = "ṁ" // anusvara
	sc.Special[0x0D83] = "ḥ" // visarga
	return sc
}()
