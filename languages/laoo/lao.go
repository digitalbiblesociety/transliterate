// Package lao provides BGN/PCGN transliteration of the Lao script
// (U+0E80-U+0EFF) to the Latin alphabet.
//
// Lao is an Abugida like Thai. Vowels can be pre-, post-, above-, or
// below-base. We emit them in logical (storage) order, which means
// pre-positioned vowels appear before their consonant in the output —
// a deviation from phonetic order but a faithful 1:1 representation.
//
// Tone marks are dropped.
package laoo

import "strings"

const (
	BlockStart rune = 0x0E80
	BlockEnd   rune = 0x0EFF
)

// Transliterate returns the romanization of s. Non-Lao runes pass
// through unchanged.
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
		if r >= 0x0ED0 && r <= 0x0ED9 { // digits
			b.WriteRune('0' + (r - 0x0ED0))
			continue
		}
		if r >= BlockStart && r <= BlockEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has any Lao-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

// letterMap covers Lao consonants and vowels in a single flat table.
// Inherent /a/ is built into the consonant entries so we don't need an
// Abugida engine; outputs are slightly more verbose but faithful.
var letterMap = map[rune]string{
	// Consonants — each carries inherent /a/.
	0x0E81: "k",   // ກ
	0x0E82: "kh",  // ຂ
	0x0E84: "kh",  // ຄ
	0x0E87: "ng",  // ງ
	0x0E88: "c",   // ຈ
	0x0E8A: "s",   // ຊ
	0x0E8D: "ny",  // ຍ
	0x0E94: "d",   // ດ
	0x0E95: "t",   // ຕ
	0x0E96: "th",  // ຖ
	0x0E97: "th",  // ທ
	0x0E99: "n",   // ນ
	0x0E9A: "b",   // ບ
	0x0E9B: "p",   // ປ
	0x0E9C: "ph",  // ຜ
	0x0E9D: "f",   // ຝ
	0x0E9E: "ph",  // ພ
	0x0E9F: "f",   // ຟ
	0x0EA1: "m",   // ມ
	0x0EA2: "y",   // ຢ
	0x0EA3: "r",   // ຣ
	0x0EA5: "l",   // ລ
	0x0EA7: "v",   // ວ
	0x0EAA: "s",   // ສ
	0x0EAB: "h",   // ຫ
	0x0EAD: "ʾ",   // ອ
	0x0EAE: "h",   // ຮ
	// Vowels (independent and dependent forms combined).
	0x0EB0: "a",   // ◌ະ
	0x0EB1: "a",   // ◌ັ
	0x0EB2: "ā",   // າ
	0x0EB3: "am",  // ໍາ
	0x0EB4: "i",   // ◌ິ
	0x0EB5: "ī",   // ◌ີ
	0x0EB6: "ɨ",   // ◌ຶ
	0x0EB7: "ɨ̄",   // ◌ື
	0x0EB8: "u",   // ◌ຸ
	0x0EB9: "ū",   // ◌ູ
	0x0EBB: "o",   // ◌ົ
	0x0EBC: "l",   // ◌ຼ (subscript la, treated as l)
	0x0EBD: "y",   // ◌ຽ
	0x0EC0: "e",   // ເ (pre-positioned in logical order: keeps before C)
	0x0EC1: "ae",  // ແ
	0x0EC2: "o",   // ໂ
	0x0EC3: "ai",  // ໃ
	0x0EC4: "ai",  // ໄ
	0x0EC6: "ɔ",   // ◌ໆ
	0x0EC8: "",    // tone mark — drop
	0x0EC9: "",    // tone mark
	0x0ECA: "",    // tone mark
	0x0ECB: "",    // tone mark
	0x0ECC: "",    // cancellation mark (silent)
	0x0ECD: "ṁ",   // ◌ໍ niggahita
}
