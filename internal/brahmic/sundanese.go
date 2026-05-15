package brahmic

// Sundanese (U+1B80..U+1BBF).
//
// Used to write the Sundanese language of West Java. Traditional script
// (Aksara Sunda); Latin has been dominant since the 19th century but the
// traditional script is in revival.
//
// Sundanese has three medial consonant signs (pamingkal -y-, panyakra -r-,
// panyiku -l-) handled through the brahmic engine's Medial map. The
// "killer" Pamaaeh (U+1BAA) acts as the script's virama.
var Sundanese = func() *Script {
	sc := &Script{
		Name:             "sundanese",
		BlockStart:       0x1B80,
		BlockEnd:         0x1BBF,
		Virama:           []rune{0x1BAA}, // ◌᮪ pamaaeh
		DigitStart:       0x1BB0,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
		Medial:           map[rune]string{},
	}

	// Independent vowels.
	sc.IndependentVowel[0x1B83] = "a"
	sc.IndependentVowel[0x1B84] = "i"
	sc.IndependentVowel[0x1B85] = "u"
	sc.IndependentVowel[0x1B86] = "ae" // letter ae
	sc.IndependentVowel[0x1B87] = "o"
	sc.IndependentVowel[0x1B88] = "é"  // letter e (with acute, "é")
	sc.IndependentVowel[0x1B89] = "eu" // letter eu

	// Consonants.
	sc.ConsonantBase[0x1B8A] = "k"
	sc.ConsonantBase[0x1B8B] = "q"
	sc.ConsonantBase[0x1B8C] = "g"
	sc.ConsonantBase[0x1B8D] = "ṅ" // nga
	sc.ConsonantBase[0x1B8E] = "c"
	sc.ConsonantBase[0x1B8F] = "j"
	sc.ConsonantBase[0x1B90] = "z"
	sc.ConsonantBase[0x1B91] = "ñ" // nya
	sc.ConsonantBase[0x1B92] = "t"
	sc.ConsonantBase[0x1B93] = "d"
	sc.ConsonantBase[0x1B94] = "n"
	sc.ConsonantBase[0x1B95] = "p"
	sc.ConsonantBase[0x1B96] = "f"
	sc.ConsonantBase[0x1B97] = "v"
	sc.ConsonantBase[0x1B98] = "b"
	sc.ConsonantBase[0x1B99] = "m"
	sc.ConsonantBase[0x1B9A] = "y"
	sc.ConsonantBase[0x1B9B] = "r"
	sc.ConsonantBase[0x1B9C] = "l"
	sc.ConsonantBase[0x1B9D] = "w"
	sc.ConsonantBase[0x1B9E] = "s"
	sc.ConsonantBase[0x1B9F] = "x"
	sc.ConsonantBase[0x1BA0] = "h"
	sc.ConsonantBase[0x1BAE] = "kh"
	sc.ConsonantBase[0x1BAF] = "sy"

	// Medial consonant signs.
	sc.Medial[0x1BA1] = "y" // pamingkal
	sc.Medial[0x1BA2] = "r" // panyakra
	sc.Medial[0x1BA3] = "l" // panyiku

	// Vowel signs.
	sc.VowelSign[0x1BA4] = "i"  // panghulu
	sc.VowelSign[0x1BA5] = "u"  // panyuku
	sc.VowelSign[0x1BA6] = "é"  // panaelaeng
	sc.VowelSign[0x1BA7] = "o"  // panolong
	sc.VowelSign[0x1BA8] = "ə"  // pamepet
	sc.VowelSign[0x1BA9] = "eu" // paneuleung

	// Special.
	sc.Special[0x1B80] = "ṁ" // panyecek (anusvara)
	sc.Special[0x1B81] = "r" // panglayar (r-final)
	sc.Special[0x1B82] = "ḥ" // pangwisad (visarga)
	return sc
}()
