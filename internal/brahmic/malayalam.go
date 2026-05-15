package brahmic

// Malayalam (U+0D00..U+0D7F).
//
// Notes:
//   - Has both short and long e/o, like Kannada/Tamil/Telugu.
//   - Chillu letters (ൻ ർ ൽ ൾ ൺ, U+0D7A..U+0D7F) are word-final consonants
//     equivalent to their consonant + virama. Mapped via ChilluLike so they
//     don't try to consume a following vowel sign.
//   - U+0D02 is anusvara; U+0D03 is visarga (standard offsets).
var Malayalam = func() *Script {
	sc := newStandardScript("malayalam", 0x0D00)
	sc.ChilluLike = map[rune]string{
		0x0D7A: "ṇ", // ൺ
		0x0D7B: "n", // ൻ
		0x0D7C: "r", // ർ
		0x0D7D: "l", // ൽ
		0x0D7E: "ḷ", // ൾ
		0x0D7F: "k", // ൿ (rare)
	}
	return sc
}()
