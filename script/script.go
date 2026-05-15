// Package script provides a script-keyed registry of transliteration
// engines, the public counterpart to the language packages in
// `languages/`. It is the entry point external pipelines should use
// when they have an ISO 15924 script tag (e.g. from a USFM `\h` block,
// a DBL metadata field, or a conf.toml) and want the matching
// Latin-script output without manually wiring each language package.
//
// Each Engine knows:
//
//   - Name           — its ISO 15924 four-letter tag, e.g. "Thai", "Grek".
//   - Contains       — whether a given rune is part of this script's
//     block. Used by Detect to pick the dominant
//     script of arbitrary text.
//   - Transliterate  — convert source-script text to Latin.
//   - Split          — break a string into phrase-level fragments
//     suitable for forced alignment (Aeneas et al).
//     Script-specific where the writing system has its
//     own pause/segmentation conventions (currently
//     only Thai); otherwise the generic punctuation+
//     vertical-whitespace splitter from `phrases`.
//
// The CLI in cmd/translit and external tools (e.g. audio-sync) share
// this registry, so an engine added here is immediately picked up by
// every dispatcher without further wiring.
package script

import (
	"github.com/digitalbiblesociety/transliterate/languages/arab"
	"github.com/digitalbiblesociety/transliterate/languages/armn"
	"github.com/digitalbiblesociety/transliterate/languages/bali"
	"github.com/digitalbiblesociety/transliterate/languages/batk"
	"github.com/digitalbiblesociety/transliterate/languages/beng"
	"github.com/digitalbiblesociety/transliterate/languages/brah"
	"github.com/digitalbiblesociety/transliterate/languages/bugi"
	"github.com/digitalbiblesociety/transliterate/languages/cans"
	"github.com/digitalbiblesociety/transliterate/languages/cher"
	"github.com/digitalbiblesociety/transliterate/languages/cyrl"
	"github.com/digitalbiblesociety/transliterate/languages/deva"
	"github.com/digitalbiblesociety/transliterate/languages/ethi"
	"github.com/digitalbiblesociety/transliterate/languages/geor"
	"github.com/digitalbiblesociety/transliterate/languages/grek"
	"github.com/digitalbiblesociety/transliterate/languages/gujr"
	"github.com/digitalbiblesociety/transliterate/languages/guru"
	"github.com/digitalbiblesociety/transliterate/languages/hang"
	"github.com/digitalbiblesociety/transliterate/languages/hani"
	"github.com/digitalbiblesociety/transliterate/languages/hebr"
	"github.com/digitalbiblesociety/transliterate/languages/java"
	"github.com/digitalbiblesociety/transliterate/languages/jpan"
	"github.com/digitalbiblesociety/transliterate/languages/khmr"
	"github.com/digitalbiblesociety/transliterate/languages/knda"
	"github.com/digitalbiblesociety/transliterate/languages/lana"
	"github.com/digitalbiblesociety/transliterate/languages/laoo"
	"github.com/digitalbiblesociety/transliterate/languages/mlym"
	"github.com/digitalbiblesociety/transliterate/languages/modi"
	"github.com/digitalbiblesociety/transliterate/languages/mymr"
	"github.com/digitalbiblesociety/transliterate/languages/newa"
	"github.com/digitalbiblesociety/transliterate/languages/orya"
	"github.com/digitalbiblesociety/transliterate/languages/shrd"
	"github.com/digitalbiblesociety/transliterate/languages/sinh"
	"github.com/digitalbiblesociety/transliterate/languages/sund"
	"github.com/digitalbiblesociety/transliterate/languages/syrc"
	"github.com/digitalbiblesociety/transliterate/languages/taml"
	"github.com/digitalbiblesociety/transliterate/languages/telu"
	"github.com/digitalbiblesociety/transliterate/languages/thai"
	"github.com/digitalbiblesociety/transliterate/languages/tibt"
	"github.com/digitalbiblesociety/transliterate/languages/tirh"
	"github.com/digitalbiblesociety/transliterate/languages/yueh"
	"github.com/digitalbiblesociety/transliterate/phrases"
)

// Engine bundles everything a caller needs to recognise and process
// one writing system. All fields are non-nil for every registered
// engine.
type Engine struct {
	Name          string             // ISO 15924 four-letter code
	Contains      func(r rune) bool  // does r belong to this script's block?
	Transliterate func(string) string
	Split         func(string) []string
}

// engines is the single source of truth for which scripts the
// transliterate module supports. Order is significant for Detect:
// scripts that share a Unicode block with another (Hani/Yueh,
// Jpan/Hani) must be listed in the order Detect should prefer.
var engines = []Engine{
	// Brahmic family (north Indic).
	{Name: "Deva", Contains: inRange(deva.BlockStart, deva.BlockEnd), Transliterate: deva.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Beng", Contains: inRange(beng.BlockStart, beng.BlockEnd), Transliterate: beng.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Guru", Contains: inRange(guru.BlockStart, guru.BlockEnd), Transliterate: guru.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Gujr", Contains: inRange(gujr.BlockStart, gujr.BlockEnd), Transliterate: gujr.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Orya", Contains: inRange(orya.BlockStart, orya.BlockEnd), Transliterate: orya.Transliterate, Split: phrases.SplitGeneric},

	// Brahmic family (south Indic).
	{Name: "Taml", Contains: inRange(taml.BlockStart, taml.BlockEnd), Transliterate: taml.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Telu", Contains: inRange(telu.BlockStart, telu.BlockEnd), Transliterate: telu.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Knda", Contains: inRange(knda.BlockStart, knda.BlockEnd), Transliterate: knda.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Mlym", Contains: inRange(mlym.BlockStart, mlym.BlockEnd), Transliterate: mlym.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Sinh", Contains: inRange(sinh.BlockStart, sinh.BlockEnd), Transliterate: sinh.Transliterate, Split: phrases.SplitGeneric},

	// Brahmic family (Southeast Asian Indic-derived).
	{Name: "Java", Contains: inRange(java.BlockStart, java.BlockEnd), Transliterate: java.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Sund", Contains: inRange(sund.BlockStart, sund.BlockEnd), Transliterate: sund.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Bali", Contains: inRange(bali.BlockStart, bali.BlockEnd), Transliterate: bali.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Bugi", Contains: inRange(bugi.BlockStart, bugi.BlockEnd), Transliterate: bugi.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Batk", Contains: inRange(batk.BlockStart, batk.BlockEnd), Transliterate: batk.Transliterate, Split: phrases.SplitGeneric},

	// Aksharamukha-derived scripts.
	{Name: "Brah", Contains: inRange(brah.BlockStart, brah.BlockEnd), Transliterate: brah.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Shrd", Contains: inRange(shrd.BlockStart, shrd.BlockEnd), Transliterate: shrd.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Modi", Contains: inRange(modi.BlockStart, modi.BlockEnd), Transliterate: modi.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Tirh", Contains: inRange(tirh.BlockStart, tirh.BlockEnd), Transliterate: tirh.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Newa", Contains: inRange(newa.BlockStart, newa.BlockEnd), Transliterate: newa.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Lana", Contains: inRange(lana.BlockStart, lana.BlockEnd), Transliterate: lana.Transliterate, Split: phrases.SplitGeneric},

	// Other Asian scripts. Thai/Lao/Khmer/Mymr have no inter-word
	// spaces in writing; Thai and Tibetan ship custom splitters that
	// know their native sentence marks (pause marks for Thai, shad for
	// Tibetan). The others fall back to SplitGeneric until they grow
	// equivalents.
	{Name: "Mymr", Contains: inRange(mymr.BlockStart, mymr.BlockEnd), Transliterate: mymr.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Khmr", Contains: inRange(khmr.BlockStart, khmr.BlockEnd), Transliterate: khmr.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Laoo", Contains: inRange(laoo.BlockStart, laoo.BlockEnd), Transliterate: laoo.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Thai", Contains: inRange(thai.BlockStart, thai.BlockEnd), Transliterate: thai.Transliterate, Split: thai.SplitPhrases},
	{Name: "Tibt", Contains: inRange(tibt.BlockStart, tibt.BlockEnd), Transliterate: tibt.Transliterate, Split: tibt.SplitPhrases},
	{Name: "Hang", Contains: inRange(hang.BlockStart, hang.BlockEnd), Transliterate: hang.Transliterate, Split: phrases.SplitGeneric},

	// Japanese before Chinese so kana wins on mixed kana+kanji input;
	// pure-Han text still routes to Hani.
	{
		Name: "Jpan",
		Contains: func(r rune) bool {
			return (r >= jpan.HiraStart && r <= jpan.HiraEnd) || (r >= jpan.KataStart && r <= jpan.KataEnd)
		},
		Transliterate: jpan.Transliterate,
		Split:         phrases.SplitGeneric,
	},
	{Name: "Hani", Contains: inRange(hani.BlockStart, hani.BlockEnd), Transliterate: hani.Transliterate, Split: phrases.SplitGeneric},
	// Yueh (Cantonese, Jyutping) shares the Han block. Auto-detect won't
	// pick it because Hani is listed first and Detect stops at the first
	// matching engine. Callers opt in via ByName("Yueh").
	{Name: "Yueh", Contains: inRange(yueh.BlockStart, yueh.BlockEnd), Transliterate: yueh.Transliterate, Split: phrases.SplitGeneric},

	// Other.
	{Name: "Armn", Contains: inRange(armn.BlockStart, armn.BlockEnd), Transliterate: armn.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Cans", Contains: inRange(cans.BlockStart, cans.BlockEnd), Transliterate: cans.Transliterate, Split: phrases.SplitGeneric},
	{
		Name: "Cher",
		Contains: func(r rune) bool {
			return (r >= cher.BlockStart && r <= cher.BlockEnd) || (r >= cher.LowerBlockStart && r <= cher.LowerBlockEnd)
		},
		Transliterate: cher.Transliterate,
		Split:         phrases.SplitGeneric,
	},
	{Name: "Cyrl", Contains: inRange(cyrl.BlockStart, cyrl.BlockEnd), Transliterate: cyrl.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Ethi", Contains: inRange(ethi.BlockStart, ethi.BlockEnd), Transliterate: ethi.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Geor", Contains: inRange(geor.BlockStart, geor.BlockEnd), Transliterate: geor.Transliterate, Split: phrases.SplitGeneric},
	{
		Name: "Grek",
		Contains: func(r rune) bool {
			return (r >= grek.BlockStart && r <= grek.BlockEnd) || (r >= grek.ExtStart && r <= grek.ExtEnd)
		},
		Transliterate: grek.Transliterate,
		Split:         phrases.SplitGeneric,
	},
	{Name: "Hebr", Contains: inRange(hebr.BlockStart, hebr.BlockEnd), Transliterate: hebr.Transliterate, Split: phrases.SplitGeneric},
	{Name: "Syrc", Contains: inRange(syrc.BlockStart, syrc.BlockEnd), Transliterate: syrc.Transliterate, Split: phrases.SplitGeneric},

	// Arabic. The default transliterator is ANETAC dictionary lookup
	// with character-rule fallback — best for named entities in
	// otherwise-Latin text. For continuous vocalised Arabic, callers
	// should reach for arab.TransliterateTashkeel directly; that
	// engine isn't exposed here because it's a different intent.
	{
		Name:          "Arab",
		Contains:      func(r rune) bool { return (r >= 0x0600 && r <= 0x06FF) || (r >= 0x0750 && r <= 0x077F) },
		Transliterate: arabicDictRules,
		Split:         phrases.SplitGeneric,
	},
}

func inRange(lo, hi rune) func(r rune) bool {
	return func(r rune) bool { return r >= lo && r <= hi }
}

// arabicDictRules wraps arab.Transliterate's (string, []string, error)
// signature into the plain func(string) string the registry uses.
// Errors degrade to passing the input through unchanged — alignment
// pipelines prefer "best effort" to a hard failure.
func arabicDictRules(s string) string {
	out, _, err := arab.Transliterate(s)
	if err != nil {
		return s
	}
	return out
}

// Mode registrations live here (not in each language package) to
// avoid a circular import — language packages would otherwise need to
// import "script" to self-register.
func init() {
	RegisterMode("Arab", "tashkeel",
		"Arabic: vocalised-text engine (uses diacritics) instead of the default ANETAC name-dictionary lookup.",
		arab.TransliterateTashkeel)
	RegisterMode("Hani", "atonal",
		"Mandarin: strip Pinyin tone-mark diacritics from the default tonal output.",
		hani.TransliterateAtonal)
	RegisterMode("Yueh", "atonal",
		"Cantonese: strip Jyutping tone digits (1–6) from the default tonal output.",
		yueh.TransliterateAtonal)
	RegisterMode("Tibt", "phonetic",
		"Tibetan: THL Simplified Phonetic (approximate Lhasa pronunciation) instead of the default Wylie.",
		tibt.TransliteratePhonetic)
}

// ByName returns the engine whose Name matches the given ISO 15924
// code, case-insensitively. Returns nil when no such engine is
// registered.
func ByName(name string) *Engine {
	for i := range engines {
		if equalFold(engines[i].Name, name) {
			return &engines[i]
		}
	}
	return nil
}

// Detect counts how many runes of s belong to each registered engine
// and returns the winner. Returns nil when no engine's block matches
// any rune in s — that covers pure-Latin input plus inputs from
// scripts the module doesn't yet support.
func Detect(s string) *Engine {
	counts := make([]int, len(engines))
	for _, r := range s {
		for i := range engines {
			if engines[i].Contains(r) {
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
	return &engines[best]
}

// Names returns the ISO 15924 codes of every registered engine, in
// registration order. Useful for CLI help text and flag validation.
func Names() []string {
	out := make([]string, len(engines))
	for i, e := range engines {
		out[i] = e.Name
	}
	return out
}

// Engines returns a copy of the engine slice. Useful for callers that
// want to fold the registered set into a larger dispatch table
// (e.g. add CLI-level overrides such as Arabic tashkeel or Han atonal
// modes) without rebuilding the table by hand.
func Engines() []Engine {
	out := make([]Engine, len(engines))
	copy(out, engines)
	return out
}

// equalFold is strings.EqualFold inlined to keep this file dependency-
// free; the engine names are pure ASCII, so a byte-level fold suffices.
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
