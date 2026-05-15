package brahmic

// Kannada (U+0C80..U+0CFF).
//
// Standard layout. Already validated against real KANBIB/KANERV samples.
var Kannada = func() *Script {
	sc := newStandardScript("kannada", 0x0C80)
	// Kannada has ೞ (U+0CDE), the archaic LLLA, mapped to ḻ.
	sc.ConsonantBase[0x0CDE] = "ḻ"
	// Kannada has no 0x0CB4 (the standard ḻ slot is empty in this block).
	delete(sc.ConsonantBase, 0x0CB4)
	return sc
}()
