package arab

import (
	"strings"
	"unicode"
)

// normalizeArabic folds variant letter forms and strips tashkeel so that
// dictionary lookups are robust to common spelling variation.
func normalizeArabic(s string) string {
	s = strings.TrimSpace(s)
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		// strip tashkeel (diacritics) and tatweel
		case 'ً', 'ٌ', 'ٍ', 'َ', 'ُ',
			'ِ', 'ّ', 'ْ', 'ٓ', 'ٔ', 'ٕ',
			'ٰ', 'ـ':
			continue
		// alef variants → bare alef
		case 'أ', 'إ', 'آ', 'ٱ':
			b.WriteRune('ا')
		// alef maqsura → ya
		case 'ى':
			b.WriteRune('ي')
		// ta marbuta → ha (matches ANETAC convention for many names)
		case 'ة':
			b.WriteRune('ه')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// digraphs are multi-rune Arabic sequences that map to a single Latin string.
// Order matters: applied longest-first by ApplyRules.
var digraphs = map[string]string{
	"اا": "aa",
	"يي": "iy",
	"وو": "uw",
}

// charMap is the single-rune fallback for the rule-based transliterator.
// Values reflect common ANETAC-style romanization of names rather than
// strict Buckwalter or ISO 233.
var charMap = map[rune]string{
	'ا': "a", 'ب': "b", 'ت': "t", 'ث': "th",
	'ج': "j", 'ح': "h", 'خ': "kh", 'د': "d",
	'ذ': "dh", 'ر': "r", 'ز': "z", 'س': "s",
	'ش': "sh", 'ص': "s", 'ض': "d", 'ط': "t",
	'ظ': "z", 'ع': "a", 'غ': "gh", 'ف': "f",
	'ق': "q", 'ك': "k", 'ل': "l", 'م': "m",
	'ن': "n", 'ه': "h", 'و': "w", 'ي': "y",
	'ء': "", 'ؤ': "u", 'ئ': "i",
	// latin / digits pass through; handled below.
}

// ApplyRules transliterates an Arabic string using char/digraph rules only.
// The result is title-cased on the first letter to match name conventions.
func ApplyRules(arabic string) string {
	s := normalizeArabic(arabic)
	var b strings.Builder
	b.Grow(len(s))

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if i+1 < len(runes) {
			pair := string(runes[i : i+2])
			if v, ok := digraphs[pair]; ok {
				b.WriteString(v)
				i++
				continue
			}
		}
		r := runes[i]
		if v, ok := charMap[r]; ok {
			b.WriteString(v)
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '\'' {
			b.WriteRune(r)
		}
	}
	return titleName(b.String())
}

// titleName capitalizes the first letter of each whitespace- or hyphen-
// separated word.
func titleName(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	capNext := true
	for _, r := range s {
		if r == ' ' || r == '-' {
			capNext = true
			b.WriteRune(r)
			continue
		}
		if capNext {
			b.WriteRune(unicode.ToUpper(r))
			capNext = false
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
