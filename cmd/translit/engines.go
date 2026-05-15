package main

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
	"github.com/digitalbiblesociety/transliterate/languages/laoo"
	"github.com/digitalbiblesociety/transliterate/languages/mlym"
	"github.com/digitalbiblesociety/transliterate/languages/mymr"
	"github.com/digitalbiblesociety/transliterate/languages/orya"
	"github.com/digitalbiblesociety/transliterate/languages/shrd"
	"github.com/digitalbiblesociety/transliterate/languages/sinh"
	"github.com/digitalbiblesociety/transliterate/languages/sund"
	"github.com/digitalbiblesociety/transliterate/languages/syrc"
	"github.com/digitalbiblesociety/transliterate/languages/taml"
	"github.com/digitalbiblesociety/transliterate/languages/telu"
	"github.com/digitalbiblesociety/transliterate/languages/thai"
	"github.com/digitalbiblesociety/transliterate/languages/tibt"
	"github.com/digitalbiblesociety/transliterate/languages/yueh"
)

// engine bundles everything needed to recognize and transliterate one
// writing system. Auto-detection picks the engine whose inBlock matches
// the most runes of the input.
//
// Names are ISO 15924 four-letter script codes; engineByName performs a
// case-insensitive lookup so users can type "grek" or "Grek".
type engine struct {
	name          string
	inBlock       func(r rune) bool
	transliterate func(string) string

	// tashkeel is an alternate transliteration function selected by the
	// CLI's -tashkeel flag. Set only for Arab; nil elsewhere.
	tashkeel func(string) string

	// atonal is an alternate transliteration function selected by the
	// CLI's -notones flag. Set for Hani / Yueh; nil elsewhere.
	atonal func(string) string
}

var engines = []engine{
	// Brahmic family (north Indic).
	{name: "Deva", inBlock: inRange(deva.BlockStart, deva.BlockEnd), transliterate: deva.Transliterate},
	{name: "Beng", inBlock: inRange(beng.BlockStart, beng.BlockEnd), transliterate: beng.Transliterate},
	{name: "Guru", inBlock: inRange(guru.BlockStart, guru.BlockEnd), transliterate: guru.Transliterate},
	{name: "Gujr", inBlock: inRange(gujr.BlockStart, gujr.BlockEnd), transliterate: gujr.Transliterate},
	{name: "Orya", inBlock: inRange(orya.BlockStart, orya.BlockEnd), transliterate: orya.Transliterate},

	// Brahmic family (south Indic).
	{name: "Taml", inBlock: inRange(taml.BlockStart, taml.BlockEnd), transliterate: taml.Transliterate},
	{name: "Telu", inBlock: inRange(telu.BlockStart, telu.BlockEnd), transliterate: telu.Transliterate},
	{name: "Knda", inBlock: inRange(knda.BlockStart, knda.BlockEnd), transliterate: knda.Transliterate},
	{name: "Mlym", inBlock: inRange(mlym.BlockStart, mlym.BlockEnd), transliterate: mlym.Transliterate},
	{name: "Sinh", inBlock: inRange(sinh.BlockStart, sinh.BlockEnd), transliterate: sinh.Transliterate},

	// Brahmic family (Southeast Asian Indic-derived).
	{name: "Java", inBlock: inRange(java.BlockStart, java.BlockEnd), transliterate: java.Transliterate},
	{name: "Sund", inBlock: inRange(sund.BlockStart, sund.BlockEnd), transliterate: sund.Transliterate},
	{name: "Bali", inBlock: inRange(bali.BlockStart, bali.BlockEnd), transliterate: bali.Transliterate},
	{name: "Bugi", inBlock: inRange(bugi.BlockStart, bugi.BlockEnd), transliterate: bugi.Transliterate},
	{name: "Batk", inBlock: inRange(batk.BlockStart, batk.BlockEnd), transliterate: batk.Transliterate},

	// Aksharamukha-derived historical Indic scripts.
	{name: "Brah", inBlock: inRange(brah.BlockStart, brah.BlockEnd), transliterate: brah.Transliterate},
	{name: "Shrd", inBlock: inRange(shrd.BlockStart, shrd.BlockEnd), transliterate: shrd.Transliterate},

	// Other Asian scripts.
	{name: "Mymr", inBlock: inRange(mymr.BlockStart, mymr.BlockEnd), transliterate: mymr.Transliterate},
	{name: "Khmr", inBlock: inRange(khmr.BlockStart, khmr.BlockEnd), transliterate: khmr.Transliterate},
	{name: "Laoo", inBlock: inRange(laoo.BlockStart, laoo.BlockEnd), transliterate: laoo.Transliterate},
	{name: "Thai", inBlock: inRange(thai.BlockStart, thai.BlockEnd), transliterate: thai.Transliterate},
	{name: "Tibt", inBlock: inRange(tibt.BlockStart, tibt.BlockEnd), transliterate: tibt.Transliterate},
	{name: "Hang", inBlock: inRange(hang.BlockStart, hang.BlockEnd), transliterate: hang.Transliterate},

	// Japanese before Chinese so kana wins on mixed kana+kanji input;
	// pure-Han text still routes to Hani.
	{name: "Jpan", inBlock: func(r rune) bool {
		return (r >= jpan.HiraStart && r <= jpan.HiraEnd) || (r >= jpan.KataStart && r <= jpan.KataEnd)
	}, transliterate: jpan.Transliterate},
	{name: "Hani", inBlock: inRange(hani.BlockStart, hani.BlockEnd),
		transliterate: hani.Transliterate, atonal: hani.TransliterateAtonal},
	// Yueh (Cantonese, Jyutping) shares the Han block. Auto-detect won't
	// pick it because Hani is listed first and `detect` stops at the first
	// matching engine; users opt in via `-script Yueh`.
	{name: "Yueh", inBlock: inRange(yueh.BlockStart, yueh.BlockEnd),
		transliterate: yueh.Transliterate, atonal: yueh.TransliterateAtonal},

	// Other.
	{name: "Armn", inBlock: inRange(armn.BlockStart, armn.BlockEnd), transliterate: armn.Transliterate},
	{name: "Cans", inBlock: inRange(cans.BlockStart, cans.BlockEnd), transliterate: cans.Transliterate},
	{name: "Cher", inBlock: func(r rune) bool {
		return (r >= cher.BlockStart && r <= cher.BlockEnd) || (r >= cher.LowerBlockStart && r <= cher.LowerBlockEnd)
	}, transliterate: cher.Transliterate},
	{name: "Cyrl", inBlock: inRange(cyrl.BlockStart, cyrl.BlockEnd), transliterate: cyrl.Transliterate},
	{name: "Ethi", inBlock: inRange(ethi.BlockStart, ethi.BlockEnd), transliterate: ethi.Transliterate},
	{name: "Geor", inBlock: inRange(geor.BlockStart, geor.BlockEnd), transliterate: geor.Transliterate},
	{name: "Grek", inBlock: func(r rune) bool {
		return (r >= grek.BlockStart && r <= grek.BlockEnd) || (r >= grek.ExtStart && r <= grek.ExtEnd)
	}, transliterate: grek.Transliterate},
	{name: "Hebr", inBlock: inRange(hebr.BlockStart, hebr.BlockEnd), transliterate: hebr.Transliterate},
	{name: "Syrc", inBlock: inRange(syrc.BlockStart, syrc.BlockEnd), transliterate: syrc.Transliterate},

	// Arabic has two flavors selected by the -tashkeel flag.
	{
		name:          "Arab",
		inBlock:       func(r rune) bool { return (r >= 0x0600 && r <= 0x06FF) || (r >= 0x0750 && r <= 0x077F) },
		transliterate: arabicDictRules,
		tashkeel:      arab.TransliterateTashkeel,
	},
}

func inRange(lo, hi rune) func(r rune) bool {
	return func(r rune) bool { return r >= lo && r <= hi }
}

// arabicDictRules is the default Arabic engine: ANETAC dictionary lookup
// with character-rule fallback. Best for named entities. For continuous
// vocalized text, pass -tashkeel.
func arabicDictRules(s string) string {
	out, _, err := arab.Transliterate(s)
	if err != nil {
		return s
	}
	return out
}

// detect picks the engine matching the most runes of s. Returns nil if
// no engine's block has any runes in s.
func detect(s string) *engine {
	counts := make([]int, len(engines))
	for _, r := range s {
		for i := range engines {
			if engines[i].inBlock(r) {
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

// engineByName looks up an engine by ISO 15924 code, case-insensitively.
// Returns nil if not found.
func engineByName(name string) *engine {
	for i := range engines {
		if equalFold(engines[i].name, name) {
			return &engines[i]
		}
	}
	return nil
}

// equalFold is strings.EqualFold inlined to keep this file dependency-free.
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

// engineNames returns the ISO 15924 codes of every registered engine,
// useful for help text and validation.
func engineNames() []string {
	names := make([]string, len(engines))
	for i, e := range engines {
		names[i] = e.name
	}
	return names
}
