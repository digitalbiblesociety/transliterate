package arab

import (
	"strings"
	"unicode"
)

// USFMLine transliterates the Arabic content of a USFM line while preserving
// markers (tokens beginning with a backslash) and their argument structure
// verbatim. Tokens are split on whitespace; whitespace runs are preserved.
//
// Behavior per token:
//   - starts with '\' → emitted unchanged (USFM marker, e.g. \v, \f*, \fr)
//   - contains Arabic-script runes → leading/trailing punctuation is
//     preserved, the Arabic core is sent through t.Transliterate
//   - otherwise → emitted unchanged (Latin punctuation, Western digits, etc.)
//
// Arabic-Indic digits (٠-٩) inside non-marker tokens are converted to
// Western digits.
func (t *Transliterator) USFMLine(line string) string {
	if line == "" {
		return line
	}
	var b strings.Builder
	b.Grow(len(line) * 2)

	i := 0
	runes := []rune(line)
	for i < len(runes) {
		if unicode.IsSpace(runes[i]) {
			start := i
			for i < len(runes) && unicode.IsSpace(runes[i]) {
				i++
			}
			b.WriteString(string(runes[start:i]))
			continue
		}
		start := i
		for i < len(runes) && !unicode.IsSpace(runes[i]) {
			i++
		}
		tok := string(runes[start:i])
		b.WriteString(t.transformToken(tok))
	}
	return b.String()
}

func (t *Transliterator) transformToken(tok string) string {
	if strings.HasPrefix(tok, "\\") {
		return tok // USFM marker, untouched
	}
	if !containsArabic(tok) {
		return convertArabicDigits(tok)
	}

	// Split into [leading-punct] [arabic-core] [trailing-punct], where the
	// core is the longest run from the first to the last Arabic letter.
	rs := []rune(tok)
	firstAr, lastAr := -1, -1
	for i, r := range rs {
		if isArabicLetter(r) {
			if firstAr == -1 {
				firstAr = i
			}
			lastAr = i
		}
	}
	if firstAr == -1 {
		return convertArabicDigits(tok)
	}
	lead := string(rs[:firstAr])
	core := string(rs[firstAr : lastAr+1])
	trail := string(rs[lastAr+1:])

	out, _ := t.Transliterate(core)
	return convertArabicDigits(lead) + out + convertArabicDigits(trail)
}

func isArabicLetter(r rune) bool {
	if (r >= 0x0600 && r <= 0x06FF) || (r >= 0x0750 && r <= 0x077F) {
		// exclude digits and diacritics from "letter" classification so
		// punctuation around a number doesn't get folded into the core
		if r >= 0x0660 && r <= 0x0669 {
			return false
		}
		return unicode.IsLetter(r) || (r >= 0x0600 && r <= 0x06FF)
	}
	return false
}

func convertArabicDigits(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 0x0660 && r <= 0x0669: // Arabic-Indic 0-9
			b.WriteRune('0' + (r - 0x0660))
		case r >= 0x06F0 && r <= 0x06F9: // Extended Arabic-Indic 0-9
			b.WriteRune('0' + (r - 0x06F0))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
