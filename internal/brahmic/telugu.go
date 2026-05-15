package brahmic

// Telugu (U+0C00..U+0C7F).
//
// Highly parallel layout with Kannada — same short/long e and o
// distinction native to South Indian scripts.
var Telugu = func() *Script {
	sc := newStandardScript("telugu", 0x0C00)
	// Telugu lacks the archaic ṟ at +0x31 (it does have ఱ at +0x31 actually
	// — keep it). ḻ at +0x34 not in modern use; remove.
	delete(sc.ConsonantBase, 0x0C34)
	// Telugu length marks U+0C55/0C56 — drop silently.
	sc.Special[0x0C55] = ""
	sc.Special[0x0C56] = ""
	return sc
}()
