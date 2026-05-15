package brahmic

// Javanese (U+A980..U+A9DF).
//
// Used to write the Javanese language and Kawi (Old Javanese). Latin has
// been the everyday script since the 19th-century Dutch colonial era, but
// Aksara Jawa is taught in schools and used for signage in Yogyakarta and
// Central Java.
//
// Javanese has two medial consonant signs handled through the brahmic
// engine's Medial map:
//   - Pengkal (U+A9BE) attaches a -y- after the head consonant.
//   - Cakra (U+A9BF) attaches a -r- after the head consonant.
//
// The script encodes "murda" (noble) and "mahaprana" (aspirated) variant
// forms of several consonants as distinct codepoints. For ISO 15919-style
// romanization we collapse murda variants to the same Latin base; the
// mahaprana set gets the aspirated form (kh, gh, etc.) per Sanskrit
// convention.
var Javanese = func() *Script {
	sc := &Script{
		Name:             "javanese",
		BlockStart:       0xA980,
		BlockEnd:         0xA9DF,
		Virama:           []rune{0xA9C0}, // pangkon
		DigitStart:       0xA9D0,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
		Medial:           map[rune]string{},
	}

	// Independent vowels.
	sc.IndependentVowel[0xA984] = "a"
	sc.IndependentVowel[0xA985] = "i" // i kawi
	sc.IndependentVowel[0xA986] = "i"
	sc.IndependentVowel[0xA987] = "ī"
	sc.IndependentVowel[0xA988] = "u"
	sc.IndependentVowel[0xA989] = "r̥" // pa cerek
	sc.IndependentVowel[0xA98A] = "l̥" // nga lelet
	sc.IndependentVowel[0xA98B] = "l̥̄"
	sc.IndependentVowel[0xA98C] = "e"
	sc.IndependentVowel[0xA98D] = "ai"
	sc.IndependentVowel[0xA98E] = "o"

	// Consonants (with murda variants collapsing to the same Latin base).
	sc.ConsonantBase[0xA98F] = "k"
	sc.ConsonantBase[0xA990] = "k"  // ka sasak
	sc.ConsonantBase[0xA991] = "k"  // ka murda
	sc.ConsonantBase[0xA992] = "g"
	sc.ConsonantBase[0xA993] = "g"  // ga murda
	sc.ConsonantBase[0xA994] = "ṅ"  // nga
	sc.ConsonantBase[0xA995] = "c"
	sc.ConsonantBase[0xA996] = "c"  // ca murda
	sc.ConsonantBase[0xA997] = "j"
	sc.ConsonantBase[0xA998] = "ñ"  // nya murda
	sc.ConsonantBase[0xA999] = "jh" // ja mahaprana
	sc.ConsonantBase[0xA99A] = "ñ"  // nya
	sc.ConsonantBase[0xA99B] = "ṭ"
	sc.ConsonantBase[0xA99C] = "ṭh" // tta mahaprana
	sc.ConsonantBase[0xA99D] = "ḍ"
	sc.ConsonantBase[0xA99E] = "ḍh" // dda mahaprana
	sc.ConsonantBase[0xA99F] = "ṇ"  // na murda
	sc.ConsonantBase[0xA9A0] = "t"
	sc.ConsonantBase[0xA9A1] = "t"  // ta murda
	sc.ConsonantBase[0xA9A2] = "d"
	sc.ConsonantBase[0xA9A3] = "dh" // da mahaprana
	sc.ConsonantBase[0xA9A4] = "n"
	sc.ConsonantBase[0xA9A5] = "p"
	sc.ConsonantBase[0xA9A6] = "p" // pa murda
	sc.ConsonantBase[0xA9A7] = "b"
	sc.ConsonantBase[0xA9A8] = "b" // ba murda
	sc.ConsonantBase[0xA9A9] = "m"
	sc.ConsonantBase[0xA9AA] = "y"
	sc.ConsonantBase[0xA9AB] = "r"
	sc.ConsonantBase[0xA9AC] = "r" // ra agung
	sc.ConsonantBase[0xA9AD] = "l"
	sc.ConsonantBase[0xA9AE] = "w"
	sc.ConsonantBase[0xA9AF] = "ś"  // sa murda
	sc.ConsonantBase[0xA9B0] = "ṣ"  // sa mahaprana
	sc.ConsonantBase[0xA9B1] = "s"
	sc.ConsonantBase[0xA9B2] = "h"

	// Vowel signs (matras).
	sc.VowelSign[0xA9B4] = "ā" // tarung (lengthens the previous vowel; we approximate as ā)
	sc.VowelSign[0xA9B5] = "o" // tolong
	sc.VowelSign[0xA9B6] = "i" // wulu
	sc.VowelSign[0xA9B7] = "ī" // wulu melik
	sc.VowelSign[0xA9B8] = "u" // suku
	sc.VowelSign[0xA9B9] = "ū" // suku mendut
	sc.VowelSign[0xA9BA] = "e" // taling
	sc.VowelSign[0xA9BB] = "ai" // dirga mure
	sc.VowelSign[0xA9BC] = "ə" // pepet
	sc.VowelSign[0xA9BD] = "r̥" // keret

	// Medial consonant signs (consumed before vowel sign / virama).
	sc.Medial[0xA9BE] = "y" // pengkal
	sc.Medial[0xA9BF] = "r" // cakra

	// Special markers.
	sc.Special[0xA980] = "ṁ" // panyangga (anusvara-like)
	sc.Special[0xA981] = "ṁ" // cecak (ng-final)
	sc.Special[0xA982] = "r" // layar (r-final)
	sc.Special[0xA983] = "ḥ" // wignyan (h-final)
	return sc
}()
