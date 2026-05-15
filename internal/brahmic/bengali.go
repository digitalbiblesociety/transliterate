package brahmic

// Bengali (U+0980..U+09FF): Bengali, Assamese, Manipuri (Bishnupriya).
//
// Notes:
//   - Bengali script has no short e/o at +0x0E/+0x12 (parallel to
//     Devanagari).
//   - Khanda Ta (ৎ, U+09CE) is a final-consonant form equivalent to ত +
//     virama. Mapped via ChilluLike so it doesn't try to consume a
//     following vowel sign.
//   - The Bengali letter ৱ (U+09F1, "wa") is Assamese; mapped to "wa".
//   - Bengali numerals U+09E6..U+09EF are at the standard digit offset.
//   - Bengali "Ya" (য, U+09AF) is "ya" but the Bengali "Ya with hook" (য়,
//     U+09DF) is "ẏa" in ISO 15919.
var Bengali = func() *Script {
	sc := newStandardScript("bengali", 0x0980)
	delete(sc.IndependentVowel, 0x098E)
	delete(sc.IndependentVowel, 0x0992)
	delete(sc.VowelSign, 0x09C6)
	delete(sc.VowelSign, 0x09CA)
	sc.IndependentVowel[0x098F] = "e"
	sc.IndependentVowel[0x0993] = "o"
	sc.VowelSign[0x09C7] = "e"
	sc.VowelSign[0x09CB] = "o"
	// Bengali Ra is at +0x30 (র) — already standard. But Bengali also has
	// ৰ (U+09F0) for Assamese Ra and ৱ (U+09F1) for Assamese Wa.
	sc.ConsonantBase[0x09F0] = "r"  // Assamese ra
	sc.ConsonantBase[0x09F1] = "w"  // Assamese wa
	sc.ConsonantBase[0x09DC] = "ṛ"  // ড় (ra-with-dot)
	sc.ConsonantBase[0x09DD] = "ṛh" // ঢ় (rha-with-dot)
	sc.ConsonantBase[0x09DF] = "ẏ"  // য় (ya with hook)
	// Bengali script does not have ṟ / ḻ / archaic LLLA at +0x31/+0x34.
	delete(sc.ConsonantBase, 0x09B1)
	delete(sc.ConsonantBase, 0x09B4)
	// Khanda Ta — final consonant; behaves as t̪ + virama.
	sc.ChilluLike = map[rune]string{
		0x09CE: "t", // ৎ
	}
	// Nukta combining (decomposed form).
	sc.Nukta = 0x09BC
	sc.NuktaCombine = map[rune]string{
		0x09A1: "ṛ",  // ড + ◌়
		0x09A2: "ṛh", // ঢ + ◌়
		0x09AF: "ẏ",  // য + ◌়
	}
	return sc
}()
