// Package armenian provides ISO 9985 transliteration of the Armenian
// alphabet (U+0530-U+058F) to the Latin alphabet. The alphabet has 39
// letters with deterministic 1:1 mappings.
package armn

import "strings"

const (
	BlockStart rune = 0x0530
	BlockEnd   rune = 0x058F
)

// Transliterate returns the ISO 9985 romanization of s. Non-Armenian
// runes pass through unchanged.
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
		if r >= BlockStart && r <= BlockEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has any Armenian-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

var letterMap = map[rune]string{
	// Uppercase.
	0x0531: "A",  // Ա
	0x0532: "B",  // Բ
	0x0533: "G",  // Գ
	0x0534: "D",  // Դ
	0x0535: "E",  // Ե
	0x0536: "Z",  // Զ
	0x0537: "Ē",  // Է
	0x0538: "Ə",  // Ը
	0x0539: "Tʿ", // Թ
	0x053A: "Ž",  // Ժ
	0x053B: "I",  // Ի
	0x053C: "L",  // Լ
	0x053D: "X",  // Խ
	0x053E: "C",  // Ծ
	0x053F: "K",  // Կ
	0x0540: "H",  // Հ
	0x0541: "J",  // Ձ
	0x0542: "Ł",  // Ղ
	0x0543: "Č",  // Ճ
	0x0544: "M",  // Մ
	0x0545: "Y",  // Յ
	0x0546: "N",  // Ն
	0x0547: "Š",  // Շ
	0x0548: "O",  // Ո
	0x0549: "Čʿ", // Չ
	0x054A: "P",  // Պ
	0x054B: "J̌",  // Ջ
	0x054C: "Ṙ",  // Ռ
	0x054D: "S",  // Ս
	0x054E: "V",  // Վ
	0x054F: "T",  // Տ
	0x0550: "R",  // Ր
	0x0551: "Cʿ", // Ց
	0x0552: "W",  // Ւ
	0x0553: "Pʿ", // Փ
	0x0554: "Kʿ", // Ք
	0x0555: "Ō",  // Օ
	0x0556: "F",  // Ֆ
	// Lowercase.
	0x0561: "a",
	0x0562: "b",
	0x0563: "g",
	0x0564: "d",
	0x0565: "e",
	0x0566: "z",
	0x0567: "ē",
	0x0568: "ə",
	0x0569: "tʿ",
	0x056A: "ž",
	0x056B: "i",
	0x056C: "l",
	0x056D: "x",
	0x056E: "c",
	0x056F: "k",
	0x0570: "h",
	0x0571: "j",
	0x0572: "ł",
	0x0573: "č",
	0x0574: "m",
	0x0575: "y",
	0x0576: "n",
	0x0577: "š",
	0x0578: "o",
	0x0579: "čʿ",
	0x057A: "p",
	0x057B: "ǰ",
	0x057C: "ṙ",
	0x057D: "s",
	0x057E: "v",
	0x057F: "t",
	0x0580: "r",
	0x0581: "cʿ",
	0x0582: "w",
	0x0583: "pʿ",
	0x0584: "kʿ",
	0x0585: "ō",
	0x0586: "f",
	0x0587: "ew", // և yev ligature
	// Punctuation.
	0x0589: ".",  // ։ full stop
	0x055D: ",",  // ՝ comma
	0x055E: "?",  // ՞ question
	0x055C: "!",  // ՜ exclamation
	0x058A: "-",  // ֊ Armenian hyphen
}
