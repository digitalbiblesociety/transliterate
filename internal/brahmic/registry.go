package brahmic

// All lists every Script in this package. Order is significant for
// auto-detection: scripts that share rune ranges with another (none in
// practice for Brahmic) would be matched in this order.
var All = []*Script{
	Devanagari,
	Bengali,
	Gurmukhi,
	Gujarati,
	Oriya,
	Tamil,
	Telugu,
	Kannada,
	Malayalam,
	Sinhala,
	Javanese,
	Sundanese,
	Balinese,
	Buginese,
	Batak,
	Brahmi,
	Sharada,
}

// Detect scans s and returns the Script whose block contains the most
// runes. Returns nil if s has no Brahmic content. Useful for picking a
// transliterator when processing a USFM file whose script isn't declared.
func Detect(s string) *Script {
	counts := make([]int, len(All))
	for _, r := range s {
		for i, sc := range All {
			if r >= sc.BlockStart && r <= sc.BlockEnd {
				counts[i]++
				break
			}
		}
	}
	best, bestN := -1, 0
	for i, n := range counts {
		if n > bestN {
			best, bestN = i, n
		}
	}
	if best < 0 {
		return nil
	}
	return All[best]
}
