package hebr

import (
	"strings"
	"unicode"
)

// cluster is the working unit of Hebrew text: one base rune (a Hebrew
// letter or pass-through character) plus the combining marks that
// follow it.
type cluster struct {
	cons    rune // base letter, or non-Hebrew rune for pass-through
	dagesh  bool // U+05BC dagesh / mappiq on this cluster
	shinSin rune // ShinDot, SinDot, or 0
	vowel   rune // primary niqqud, 0 if none
}

// isHebrew returns true if the cluster's base is one of the 27 Hebrew
// letters. Punctuation, whitespace, and out-of-block runes return
// false.
func (c cluster) isHebrew() bool { return isHebrewLetter(c.cons) }

// parseClusters walks an NFD-decomposed rune slice and groups each
// base rune with any following combining marks.
func parseClusters(rs []rune) []cluster {
	out := make([]cluster, 0, len(rs))
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		c := cluster{cons: r}
		if isHebrewLetter(r) {
			// Consume the combining-mark cluster that follows.
			for i+1 < len(rs) && unicode.Is(unicode.Mn, rs[i+1]) {
				i++
				m := rs[i]
				switch {
				case isVowelMark(m):
					// Only keep the first vowel seen — duplicate niqqud
					// is malformed input, and we want deterministic
					// output rather than a last-write-wins surprise.
					if c.vowel == 0 {
						c.vowel = m
					}
				case m == Dagesh:
					c.dagesh = true
				case m == ShinDot, m == SinDot:
					c.shinSin = m
				}
				// Other marks (meteg, rafe, cantillation) are silently
				// dropped.
			}
		}
		out = append(out, c)
	}
	return out
}

// annotation carries the per-cluster decisions made by analyse.
type annotation struct {
	skip             bool   // this cluster was absorbed as a mater
	standaloneShureq bool   // word-initial vav+dagesh acting as û on its own
	furtivePatah     bool   // emit patah before, not after, the consonant
	dageshForte      bool   // dagesh chazaq → double the consonant
	silent           bool   // shewa here is silent (don't emit "ə")
	vowelOut         string // Latin vowel to emit (may be "")
}

// emit walks the cluster stream produced by parseClusters and writes
// the SBL Latin form into b. Hebrew clusters are batched into words so
// word-level rules (divine name, furtive patah, qamatz-katan via
// maqaf) can be applied.
func emit(cs []cluster, b *strings.Builder) {
	i := 0
	n := len(cs)
	for i < n {
		if !cs[i].isHebrew() {
			b.WriteString(emitNonHebrew(cs[i]))
			i++
			continue
		}
		j := i
		for j < n && cs[j].isHebrew() {
			j++
		}
		// Pass the next non-Hebrew cluster (or zero-value sentinel) so
		// the qamatz-katan rule can see a trailing maqaf.
		var follower cluster
		if j < n {
			follower = cs[j]
		}
		emitWord(cs[i:j], follower, b)
		i = j
	}
}

// emitNonHebrew handles punctuation, whitespace, and any non-Hebrew
// rune that appears between (or inside) Hebrew words.
func emitNonHebrew(c cluster) string {
	if v, ok := punctuationLatin[c.cons]; ok {
		return v
	}
	// In-block but unmapped (rare stray combining mark, etc.) → drop.
	if c.cons >= BlockStart && c.cons <= BlockEnd {
		return ""
	}
	return string(c.cons)
}

// emitWord transliterates one Hebrew-letter word. follower is the next
// cluster after the word (or a zero cluster at end of input) and is
// used only by the qamatz-katan / maqaf heuristic.
func emitWord(cs []cluster, follower cluster, b *strings.Builder) {
	if isDivineName(cs) {
		b.WriteString("yhwh")
		return
	}
	ann := analyse(cs, follower)
	for i, c := range cs {
		a := ann[i]
		if a.skip {
			continue
		}
		if a.standaloneShureq {
			b.WriteString("û")
			continue
		}
		if a.furtivePatah {
			b.WriteString("a")
		}
		consStr := consonantString(c)
		if a.dageshForte {
			b.WriteString(consStr)
			b.WriteString(consStr)
		} else {
			b.WriteString(consStr)
		}
		if !a.furtivePatah {
			b.WriteString(a.vowelOut)
		}
	}
}

// consonantString returns the Latin form for c's consonant, taking
// shin/sin dot into account. Returns "" if the cluster's base isn't a
// known consonant (shouldn't happen for Hebrew-letter clusters but is
// safe).
func consonantString(c cluster) string {
	if c.cons == Shin {
		switch c.shinSin {
		case SinDot:
			return "ś"
		default:
			return "š"
		}
	}
	return consonantLatin[c.cons]
}

// isDivineName reports whether cs spells the Tetragrammaton (yod-he-
// vav-he), regardless of pointing.
func isDivineName(cs []cluster) bool {
	if len(cs) != 4 {
		return false
	}
	return cs[0].cons == Yod &&
		cs[1].cons == He &&
		cs[2].cons == Vav &&
		cs[3].cons == He
}

// analyse computes per-cluster decisions for one word. It runs three
// passes:
//  1. detect maters / ligatures (consume vav, yod, or he that serve as
//     vowel letters and rewrite the preceding cluster's vowel);
//  2. resolve each cluster's vowel string, including silent shewa and
//     qamatz-katan;
//  3. mark dagesh chazaq doublings and the word-final furtive patah.
func analyse(cs []cluster, follower cluster) []annotation {
	n := len(cs)
	ann := make([]annotation, n)

	// Pass 1: ligatures.
	for i := range n {
		c := cs[i]
		if i == 0 {
			// Word-initial shureq (וּ at the start of a word) carries
			// no preceding consonant, so it stands alone as "û".
			if c.cons == Vav && c.dagesh && c.vowel == 0 {
				ann[i].standaloneShureq = true
			}
			continue
		}
		prev := cs[i-1]
		if ann[i-1].skip || ann[i-1].standaloneShureq {
			// The "previous" cluster was absorbed; mater rules must
			// reach further back to be meaningful, and the simple
			// adjacent-prev case doesn't apply.
			continue
		}

		// Yod mater after hiriq / tsere / segol with no vowel of its own.
		if c.cons == Yod && c.vowel == 0 && !c.dagesh {
			switch prev.vowel {
			case Hiriq:
				ann[i-1].vowelOut = "î"
				ann[i].skip = true
				continue
			case Tsere:
				ann[i-1].vowelOut = "ê"
				ann[i].skip = true
				continue
			case Segol:
				ann[i-1].vowelOut = "ê"
				ann[i].skip = true
				continue
			}
		}

		// Vav mater carrying holam: cluster (vav, holam, no dagesh)
		// where the previous cluster has no vowel of its own collapses
		// to holam-vav → ô.
		if c.cons == Vav && c.vowel == Holam && !c.dagesh && prev.vowel == 0 {
			ann[i-1].vowelOut = "ô"
			ann[i].skip = true
			continue
		}

		// Shureq mater: cluster (vav, no vowel, dagesh) after a
		// vowel-less consonant supplies that consonant with û.
		if c.cons == Vav && c.dagesh && c.vowel == 0 && prev.vowel == 0 {
			ann[i-1].vowelOut = "û"
			ann[i].skip = true
			continue
		}

		// Qamatz-he ligature: word-final he with no dagesh/mappiq and
		// no vowel of its own, preceded by qamatz → â.
		if i == n-1 && c.cons == He && c.vowel == 0 && !c.dagesh {
			if prev.vowel == Qamats {
				ann[i-1].vowelOut = "â"
				ann[i].skip = true
				continue
			}
		}
	}

	// Pass 2: resolve vowels (qamatz-katan, silent shewa).
	for i := range n {
		if ann[i].skip || ann[i].standaloneShureq {
			continue
		}
		c := cs[i]
		// Vowel may already be set by pass 1 (ligature).
		if ann[i].vowelOut != "" {
			continue
		}
		if c.vowel == 0 {
			continue
		}
		out := vowelLatin[c.vowel]
		if c.vowel == Qamats && isQamatsQatan(cs, ann, i, follower) {
			out = "o"
		}
		if c.vowel == Sheva && !isVocalShewa(cs, ann, i) {
			out = ""
			ann[i].silent = true
		}
		ann[i].vowelOut = out
	}

	// Pass 3: dagesh forte. We need pass 2 to have finished because
	// "previous cluster has a real vowel" depends on the silent-shewa
	// decision.
	for i := range n {
		if ann[i].skip || ann[i].standaloneShureq {
			continue
		}
		c := cs[i]
		if !c.dagesh {
			continue
		}
		if isDageshForte(cs, ann, i) {
			ann[i].dageshForte = true
		}
	}

	// Pass 4: furtive patah on the final consonant.
	last := n - 1
	for last >= 0 && (ann[last].skip || ann[last].standaloneShureq) {
		last--
	}
	if last >= 0 {
		c := cs[last]
		if c.vowel == Patah && isFurtiveTarget(c) {
			ann[last].furtivePatah = true
			// The patah itself is emitted as "a" before the consonant;
			// don't double-emit it as the vowel slot.
			ann[last].vowelOut = ""
		}
	}

	return ann
}

// isFurtiveTarget reports whether c is one of the consonants that can
// host a furtive patah: ḥet, ʿayin, or he with mappiq.
func isFurtiveTarget(c cluster) bool {
	switch c.cons {
	case Het, Ayin:
		return true
	case He:
		return c.dagesh // he with mappiq is the consonantal he
	}
	return false
}

// isVocalShewa applies the heuristic that decides whether a shewa is
// vocal (emit ə) or silent (emit nothing). The classification only
// looks at left context that has already been processed, so calling it
// in pass-2 order is safe.
func isVocalShewa(cs []cluster, ann []annotation, i int) bool {
	// A hataf vowel is never a shewa; this guard is defensive.
	if cs[i].vowel != Sheva {
		return false
	}
	// A word-final shewa is always silent (closes the last syllable).
	if isLastActive(cs, ann, i) {
		return false
	}
	// Initial position in the word.
	if firstActive(ann, i) {
		return true
	}
	// On a consonant carrying dagesh forte the doubled letter ends one
	// syllable and starts the next, with the shewa belonging to the
	// new syllable → vocal.
	if cs[i].dagesh {
		// We can't read ann[i].dageshForte yet (pass 3 hasn't run), so
		// peek at the same condition directly.
		if dageshIsForte(cs, ann, i) {
			return true
		}
	}
	// Look back to the previous non-skipped cluster.
	p := prevActive(ann, i)
	if p < 0 {
		return true
	}
	pv := cs[p].vowel
	// Hataf vowels and "naturally long" vowels open a syllable → the
	// shewa that follows is vocal.
	if isLongVowel(pv) || isHatafVowel(pv) {
		return true
	}
	return false
}

// dageshIsForte mirrors isDageshForte but is callable from pass 2 (it
// only inspects features already settled by pass 1).
func dageshIsForte(cs []cluster, ann []annotation, i int) bool {
	if !cs[i].dagesh {
		return false
	}
	if firstActive(ann, i) {
		return false
	}
	p := prevActive(ann, i)
	if p < 0 {
		return false
	}
	prev := cs[p]
	// Dagesh after a consonant with no vowel of its own (i.e., a
	// silent shewa or unpointed letter) is lene, not forte.
	if prev.vowel == 0 {
		return false
	}
	if prev.vowel == Sheva && !isLongVowel(prev.vowel) {
		// Heuristic mirror of pass 2: if we'd classify the previous
		// shewa as silent, the dagesh is lene. We approximate by
		// checking ann[p].silent if set, otherwise re-evaluating.
		if ann[p].silent {
			return false
		}
	}
	return true
}

// isDageshForte is the canonical decision used in pass 3. It relies on
// pass 2 having resolved silent shewa.
func isDageshForte(cs []cluster, ann []annotation, i int) bool {
	if !cs[i].dagesh {
		return false
	}
	if firstActive(ann, i) {
		return false
	}
	// The U+05BC mark on a he is a mappiq — it only signals that the
	// he is consonantal (pronounced "h"); it never doubles.
	if cs[i].cons == He {
		return false
	}
	p := prevActive(ann, i)
	if p < 0 {
		return false
	}
	prev := cs[p]
	if prev.vowel == 0 {
		return false
	}
	if prev.vowel == Sheva && ann[p].silent {
		return false
	}
	// At word-end, dagesh on a consonant with silent shewa (or no
	// vowel at all) carries no audible doubling; we hide it.
	if i == len(cs)-1 {
		if ann[i].silent || cs[i].vowel == 0 {
			return false
		}
	}
	return true
}

// isQamatsQatan recognises a qamatz that should be read "o" rather
// than "ā". We honour two patterns beyond the explicit U+05C7:
//   - the next cluster in the same word has hataf-qamatz (the נָעֳמִי
//     case), and
//   - the word is a single-syllable closed form joined by maqaf to a
//     following word (כָּל־הָעָם).
func isQamatsQatan(cs []cluster, ann []annotation, i int, follower cluster) bool {
	// Explicit qamatz-qatan codepoint is handled by the vowel lookup
	// directly (vowelLatin[QamatsQatan] = "o"). This function only
	// runs for c.vowel == Qamats.

	// Pattern: qamatz followed by hataf-qamatz on a later cluster in
	// the same word.
	for j := i + 1; j < len(cs); j++ {
		if ann[j].skip {
			continue
		}
		if cs[j].vowel == HatafQamats {
			return true
		}
		// Stop at the next vowel-bearing cluster — the rule only
		// applies when the hataf-qamatz is immediately neighbouring
		// (allowing one intervening sheva-less consonant).
		if cs[j].vowel != 0 && cs[j].vowel != Sheva {
			break
		}
	}

	// Pattern: word's only vowel is this qamatz, followed by maqaf.
	if follower.cons == Maqaf && wordHasOnlyVowel(cs, ann, i) {
		return true
	}
	return false
}

// wordHasOnlyVowel reports whether the cluster at idx carries the only
// non-skip, non-sheva vowel in cs. Used to detect single-syllable
// closed words like כָּל.
func wordHasOnlyVowel(cs []cluster, ann []annotation, idx int) bool {
	for j, c := range cs {
		if j == idx || ann[j].skip {
			continue
		}
		if c.vowel != 0 && c.vowel != Sheva {
			return false
		}
	}
	return true
}

// firstActive reports whether i is the first non-skipped, non-
// standaloneShureq cluster in the slice — i.e. the syllable-leading
// consonant of the word.
func firstActive(ann []annotation, i int) bool {
	for j := range i {
		if !ann[j].skip && !ann[j].standaloneShureq {
			return false
		}
	}
	return true
}

// isLastActive reports whether i is the last non-skipped cluster in
// the word.
func isLastActive(cs []cluster, ann []annotation, i int) bool {
	for j := i + 1; j < len(cs); j++ {
		if !ann[j].skip && !ann[j].standaloneShureq {
			return false
		}
	}
	return true
}

// prevActive returns the index of the closest active cluster before
// i, or -1 if none exists.
func prevActive(ann []annotation, i int) int {
	for j := i - 1; j >= 0; j-- {
		if ann[j].skip || ann[j].standaloneShureq {
			continue
		}
		return j
	}
	return -1
}

// isLongVowel reports whether v counts as a "naturally long" vowel
// for the purpose of vocal-shewa classification. Matres lectionis
// (hiriq-yod, tsere-yod, etc.) are not separate runes — they're
// represented by their base vowel on the cluster, with the mater
// consumed in pass 1; the base vowel here is the one we see.
func isLongVowel(v rune) bool {
	switch v {
	case Qamats, Tsere, Holam, HolamHaser:
		return true
	}
	return false
}

// isHatafVowel reports whether v is one of the three reduced vowels.
func isHatafVowel(v rune) bool {
	switch v {
	case HatafSegol, HatafPatah, HatafQamats:
		return true
	}
	return false
}
