package brahmic

// Buginese / Lontara (U+1A00..U+1A1F).
//
// Used to write Bugis, Makassar, and other South Sulawesi languages. The
// script is an Abugida descended from Pallava via Kawi but has no virama
// (no inherent-/a/ suppression at all) — consonant clusters and final
// consonants are reader-inferred from context. We therefore emit every
// consonant with the inherent /a/ unless an explicit vowel sign overrides
// it; word-final consonants render with a trailing "a" that doesn't
// reflect actual pronunciation. This is the script's tradeoff, not a
// transliteration bug.
var Buginese = func() *Script {
	sc := &Script{
		Name:             "buginese",
		BlockStart:       0x1A00,
		BlockEnd:         0x1A1F,
		Virama:           nil, // no virama in the block
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}
	sc.ConsonantBase[0x1A00] = "k"  // ka
	sc.ConsonantBase[0x1A01] = "g"  // ga
	sc.ConsonantBase[0x1A02] = "ṅ"  // nga
	sc.ConsonantBase[0x1A03] = "ṅk" // ngka (prenasalized)
	sc.ConsonantBase[0x1A04] = "p"  // pa
	sc.ConsonantBase[0x1A05] = "b"  // ba
	sc.ConsonantBase[0x1A06] = "m"  // ma
	sc.ConsonantBase[0x1A07] = "mp" // mpa (prenasalized)
	sc.ConsonantBase[0x1A08] = "t"  // ta
	sc.ConsonantBase[0x1A09] = "d"  // da
	sc.ConsonantBase[0x1A0A] = "n"  // na
	sc.ConsonantBase[0x1A0B] = "nr" // nra (prenasalized)
	sc.ConsonantBase[0x1A0C] = "c"  // ca
	sc.ConsonantBase[0x1A0D] = "j"  // ja
	sc.ConsonantBase[0x1A0E] = "ñ"  // nya
	sc.ConsonantBase[0x1A0F] = "ñc" // nyca (prenasalized)
	sc.ConsonantBase[0x1A10] = "y"  // ya
	sc.ConsonantBase[0x1A11] = "r"  // ra
	sc.ConsonantBase[0x1A12] = "l"  // la
	sc.ConsonantBase[0x1A13] = "v"  // va
	sc.ConsonantBase[0x1A14] = "s"  // sa
	sc.ConsonantBase[0x1A16] = "h"  // ha

	sc.IndependentVowel[0x1A15] = "a" // letter A

	sc.VowelSign[0x1A17] = "i" // ◌ᨗ
	sc.VowelSign[0x1A18] = "u" // ◌ᨘ
	sc.VowelSign[0x1A19] = "e" // ◌ᨙ
	sc.VowelSign[0x1A1A] = "o" // ◌ᨚ
	sc.VowelSign[0x1A1B] = "ə" // ◌ᨛ pepet
	return sc
}()
