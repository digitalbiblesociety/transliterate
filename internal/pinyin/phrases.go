// This file is original to the transliterate project; it is not part of
// the upstream mozillazg/go-pinyin distribution. It loads the embedded
// phrase dictionary from mozillazg/phrase-pinyin-data (MIT) and exposes
// a lookup that the hani package uses to disambiguate polyphones (e.g.
// 中 → "zhōng" in 中国 but "zhòng" in 击中).

package pinyin

import (
	"embed"
	"strings"
	"sync"
)

//go:embed data/pinyin.txt data/overwrite.txt
var phraseFS embed.FS

var (
	phraseOnce sync.Once
	phraseDict map[string][]string
	maxPhrase  int
)

// MaxPhraseRunes returns the maximum rune length of any phrase in the
// embedded dictionary. Callers use it to bound longest-match scanning.
func MaxPhraseRunes() int {
	loadPhrases()
	return maxPhrase
}

// LookupPhrase returns the tone-marked pinyin syllables for s (e.g.
// ["zhōng","guó"] for "中国") if s is an exact phrase entry. The
// returned slice should be treated as read-only.
func LookupPhrase(s string) ([]string, bool) {
	loadPhrases()
	v, ok := phraseDict[s]
	return v, ok
}

func loadPhrases() {
	phraseOnce.Do(func() {
		phraseDict = make(map[string][]string, 64_000)
		// pinyin.txt first (curated primary readings), then overwrite.txt
		// (project's manual corrections). overwrite wins on conflicts.
		parsePhraseFile("data/pinyin.txt", false)
		parsePhraseFile("data/overwrite.txt", true)
	})
}

func parsePhraseFile(name string, overwrite bool) {
	b, err := phraseFS.ReadFile(name)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.Index(line, "#"); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		colon := strings.Index(line, ":")
		if colon < 0 {
			continue
		}
		phrase := strings.TrimSpace(line[:colon])
		body := strings.TrimSpace(line[colon+1:])
		if phrase == "" || body == "" {
			continue
		}
		syls := strings.Fields(body)
		if len(syls) == 0 {
			continue
		}
		// First-wins for the same phrase except in overwrite.txt where
		// later corrections supersede.
		if _, exists := phraseDict[phrase]; exists && !overwrite {
			continue
		}
		phraseDict[phrase] = syls
		if n := runeCount(phrase); n > maxPhrase {
			maxPhrase = n
		}
	}
}

func runeCount(s string) int {
	n := 0
	for range s {
		n++
	}
	return n
}
