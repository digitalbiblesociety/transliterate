// Package cyrillic provides ISO 9 (1995) transliteration of Cyrillic-
// script text to the Latin alphabet. ISO 9 is a strict 1:1 mapping
// designed to be reversible across all languages that use Cyrillic,
// including Russian, Ukrainian, Belarusian, Bulgarian, Serbian,
// Macedonian, and Mongolian.
//
// Reference: ISO 9:1995 "Information and documentation — Transliteration
// of Cyrillic characters into Latin characters — Slavic and non-Slavic
// languages."
package cyrl

import "strings"

const (
	BlockStart rune = 0x0400
	BlockEnd   rune = 0x04FF
)

// Transliterate returns the ISO 9 romanization of s. Runes outside the
// Cyrillic block pass through unchanged so the function is safe on USFM
// lines, mixed-script content, etc.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) * 2)
	for _, r := range s {
		if v, ok := letterMap[r]; ok {
			b.WriteString(v)
			continue
		}
		// In-block but unmapped (rare/historic letters, combining marks):
		// drop so we don't emit raw Cyrillic glyphs into "transliterated"
		// output.
		if r >= BlockStart && r <= BlockEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has at least one rune in the Cyrillic block.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

// letterMap encodes the ISO 9:1995 Cyrillic→Latin mapping. Where the
// standard prescribes a combining mark, we use the precomposed Latin
// character if one exists; otherwise the base letter followed by the
// combining diacritic.
var letterMap = map[rune]string{
	// Russian core alphabet.
	0x0410: "A", 0x0430: "a", // А а
	0x0411: "B", 0x0431: "b", // Б б
	0x0412: "V", 0x0432: "v", // В в
	0x0413: "G", 0x0433: "g", // Г г
	0x0414: "D", 0x0434: "d", // Д д
	0x0415: "E", 0x0435: "e", // Е е
	0x0416: "Ž", 0x0436: "ž", // Ж ж
	0x0417: "Z", 0x0437: "z", // З з
	0x0418: "I", 0x0438: "i", // И и
	0x0419: "J", 0x0439: "j", // Й й
	0x041A: "K", 0x043A: "k", // К к
	0x041B: "L", 0x043B: "l", // Л л
	0x041C: "M", 0x043C: "m", // М м
	0x041D: "N", 0x043D: "n", // Н н
	0x041E: "O", 0x043E: "o", // О о
	0x041F: "P", 0x043F: "p", // П п
	0x0420: "R", 0x0440: "r", // Р р
	0x0421: "S", 0x0441: "s", // С с
	0x0422: "T", 0x0442: "t", // Т т
	0x0423: "U", 0x0443: "u", // У у
	0x0424: "F", 0x0444: "f", // Ф ф
	0x0425: "H", 0x0445: "h", // Х х
	0x0426: "C", 0x0446: "c", // Ц ц
	0x0427: "Č", 0x0447: "č", // Ч ч
	0x0428: "Š", 0x0448: "š", // Ш ш
	0x0429: "Ŝ", 0x0449: "ŝ", // Щ щ
	0x042A: "ʺ", 0x044A: "ʺ", // Ъ ъ — hard sign
	0x042B: "Y", 0x044B: "y", // Ы ы
	0x042C: "ʹ", 0x044C: "ʹ", // Ь ь — soft sign
	0x042D: "È", 0x044D: "è", // Э э
	0x042E: "Û", 0x044E: "û", // Ю ю
	0x042F: "Â", 0x044F: "â", // Я я

	// Russian historic / extra.
	0x0401: "Ë", 0x0451: "ë", // Ё ё
	0x0462: "Ě", 0x0463: "ě", // Ѣ ѣ — yat
	0x0472: "F̀", 0x0473: "f̀", // Ѳ ѳ — fita
	0x0474: "Ỳ", 0x0475: "ỳ", // Ѵ ѵ — izhitsa
	0x0460: "O̧", 0x0461: "o̧", // Ѡ ѡ — omega

	// Ukrainian.
	0x0404: "Ê", 0x0454: "ê", // Є є
	0x0406: "Ì", 0x0456: "ì", // І і
	0x0407: "Ï", 0x0457: "ï", // Ї ї
	0x0490: "G̀", 0x0491: "g̀", // Ґ ґ

	// Belarusian.
	0x040E: "Ŭ", 0x045E: "ŭ", // Ў ў

	// Serbian / Macedonian.
	0x0402: "Đ", 0x0452: "đ", // Ђ ђ — Serbian
	0x0408: "J̌", 0x0458: "ǰ", // Ј ј
	0x0409: "L̂", 0x0459: "l̂", // Љ љ
	0x040A: "N̂", 0x045A: "n̂", // Њ њ
	0x040B: "Ć", 0x045B: "ć", // Ћ ћ — Serbian
	0x040F: "D̂", 0x045F: "d̂", // Џ џ
	0x0403: "Ǵ", 0x0453: "ǵ", // Ѓ ѓ — Macedonian
	0x040C: "Ḱ", 0x045C: "ḱ", // Ќ ќ — Macedonian
	0x0405: "Ẑ", 0x0455: "ẑ", // Ѕ ѕ — Macedonian

	// Mongolian / Tatar / Bashkir extras.
	0x04AE: "Ü", 0x04AF: "ü", // Ү ү
	0x04E8: "Ö", 0x04E9: "ö", // Ө ө
	0x0496: "Z̧", 0x0497: "z̧", // Җ җ
	0x04A2: "Ņ", 0x04A3: "ņ", // Ң ң
	0x04BA: "Ḩ", 0x04BB: "ḩ", // Һ һ
	0x049A: "Ķ", 0x049B: "ķ", // Қ қ
	0x0492: "Ǧ", 0x0493: "ǧ", // Ғ ғ
	0x04B0: "Ū", 0x04B1: "ū", // Ұ ұ
}
