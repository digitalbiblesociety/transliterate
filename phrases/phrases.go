// Package phrases provides phrase-level segmentation helpers for forced
// alignment workflows such as the Aeneas Project
// (https://github.com/readbeyond/aeneas).
//
// The typical pipeline:
//
//  1. Split source-script text into phrases (SplitGeneric or a
//     script-specific splitter like languages/thai.SplitPhrases).
//  2. Transliterate each phrase via the matching language package.
//  3. Wrap original + Latin pairs into Phrase values with stable IDs
//     via Pair.
//  4. Feed the Latin column to Aeneas; after Aeneas returns timestamps,
//     join back by ID and emit the final sync map keyed on Original.
//
// Reverse transliteration is not required — the original-script text
// travels alongside the Latin through the entire pipeline.
package phrases

import (
	"fmt"
	"strings"
)

// Phrase pairs a source-script substring with its Latin transliteration
// and a stable ID suitable for joining with forced-alignment output.
type Phrase struct {
	ID       string // stable fragment ID, e.g. "f0001"
	Original string // source-script phrase
	Latin    string // Latin transliteration
}

// SplitGeneric breaks s into phrase fragments at sentence- or clause-
// level boundaries:
//
//   - Western sentence/clause terminators: . , ; : ! ?
//   - vertical whitespace (line/paragraph separators, tabs, form feeds)
//
// Regular spaces and other horizontal whitespace (NBSP, ideographic
// space, etc.) are kept INSIDE a fragment — splitting on every space
// would produce word-level fragments that are typically too granular
// for forced alignment.
//
// Boundary characters are consumed and not included in any fragment.
// Empty / whitespace-only fragments are dropped; surviving fragments
// are trimmed of leading and trailing whitespace.
//
// SplitGeneric is the default splitter for scripts that use spaces
// between words. Scripts without inter-word spaces (Thai, Lao, Khmer,
// Burmese, CJK) should use a script-specific splitter that knows their
// native pause markers — see languages/thai.SplitPhrases as the
// reference implementation.
func SplitGeneric(s string) []string {
	var out []string
	var b strings.Builder
	flush := func() {
		if b.Len() == 0 {
			return
		}
		if frag := strings.TrimSpace(b.String()); frag != "" {
			out = append(out, frag)
		}
		b.Reset()
	}
	for _, r := range s {
		if isGenericBoundary(r) {
			flush()
			continue
		}
		b.WriteRune(r)
	}
	flush()
	return out
}

func isGenericBoundary(r rune) bool {
	switch r {
	case '.', ',', ';', ':', '!', '?':
		return true
	case '\t', '\n', '\v', '\f', '\r',
		'', // NEL
		' ', // line separator
		' ': // paragraph separator
		return true
	}
	return false
}

// Pair zips a slice of source-script phrases with their transliterations
// produced by the provided function. Each Phrase gets a stable f-prefixed
// zero-padded numeric ID (f0001, f0002, ...) for joining back with
// forced-alignment output.
//
// Example:
//
//	originals := phrases.SplitGeneric(verseText)
//	pairs := phrases.Pair(originals, grek.Transliterate)
//	// emit pairs[*].Latin to Aeneas; later join by pairs[*].ID
//	// and substitute pairs[*].Original into the final sync map.
func Pair(originals []string, transliterate func(string) string) []Phrase {
	out := make([]Phrase, len(originals))
	for i, orig := range originals {
		out[i] = Phrase{
			ID:       fmt.Sprintf("f%04d", i+1),
			Original: orig,
			Latin:    transliterate(orig),
		}
	}
	return out
}
