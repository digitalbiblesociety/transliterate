package brahmic

// Gurmukhi (U+0A00..U+0A7F): Punjabi.
//
// Notes:
//   - Gurmukhi uses Tippi (ੰ, U+0A70) and Bindi (ਂ, U+0A02) for
//     nasalization; we map both to ṁ per ISO 15919.
//   - Addak (ੱ, U+0A71) gemination is handled by emitting the next
//     consonant's base twice. Implementing that requires lookahead; for
//     now we approximate as the silent length marker (omitted).
//   - Visarga is absent in Punjabi practice; mapped anyway for completeness.
//   - Pre-composed Persian/loan consonants at U+0A59..U+0A5E.
var Gurmukhi = func() *Script {
	sc := newStandardScript("gurmukhi", 0x0A00)
	delete(sc.IndependentVowel, 0x0A0E)
	delete(sc.IndependentVowel, 0x0A12)
	delete(sc.VowelSign, 0x0A46)
	delete(sc.VowelSign, 0x0A4A)
	sc.IndependentVowel[0x0A0F] = "e"
	sc.IndependentVowel[0x0A13] = "o"
	sc.VowelSign[0x0A47] = "e"
	sc.VowelSign[0x0A4B] = "o"
	// Persian loan consonants.
	sc.ConsonantBase[0x0A59] = "k͟h"
	sc.ConsonantBase[0x0A5A] = "ġ"
	sc.ConsonantBase[0x0A5B] = "z"
	sc.ConsonantBase[0x0A5C] = "ṛ"
	sc.ConsonantBase[0x0A5E] = "f"
	// Tippi (ੰ) is also a nasalization mark.
	sc.Special[0x0A70] = "ṁ"
	// Addak (ੱ) — gemination, dropped (would need bigram pass).
	sc.Special[0x0A71] = ""
	// Gurmukhi lacks ṟ, ḻ, candrabindu, avagraha in common use.
	delete(sc.ConsonantBase, 0x0A31)
	delete(sc.ConsonantBase, 0x0A34)
	return sc
}()
