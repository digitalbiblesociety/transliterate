package hebr

import "testing"

// runCases runs Transliterate on each input and checks the result.
func runCases(t *testing.T, label string, cases map[string]string) {
	t.Helper()
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("%s: Transliterate(%q) = %q, want %q", label, in, got, want)
		}
	}
}

func TestUnpointedConsonants(t *testing.T) {
	runCases(t, "unpointed", map[string]string{
		"שלום":    "šlwm",
		"ירושלים": "yrwšlym",
		"דוד":     "dwd",
	})
}

func TestMaterLectionis(t *testing.T) {
	// Yod / vav / he serving as vowel-letters collapse with the
	// preceding vowel into the SBL ligatures.
	runCases(t, "mater", map[string]string{
		"שָׁלוֹם":     "šālôm",      // holam-vav → ô
		"בְּרֵאשִׁית": "bərēʾšît",   // hiriq-yod → î
		"אֱלֹהִים":    "ʾĕlōhîm",    // hataf-segol + hiriq-yod
		"דָּוִד":      "dāwid",      // hiriq stays short (no mater after)
		"יִשְׂרָאֵל":  "yiśrāʾēl",   // silent shewa on sin
		"קוּם":        "qûm",        // shuruq
		"עִיר":        "ʿîr",        // hiriq-yod
		"אֵין":        "ʾên",        // tsere-yod
		"עֵצָה":       "ʿēṣâ",       // qamatz-he
		"יִקְרֶה":     "yiqreh",     // segol+he stays separate
		"הָאַרְיֵה":   "hāʾaryēh",   // tsere+he stays separate
		"סוֹא":        "sôʾ",        // holam-vav with following alef
	})
}

func TestQamatsQatan(t *testing.T) {
	runCases(t, "qatan", map[string]string{
		"כָּל־הָעָם":  "kol-hāʿām",   // qamatz before maqaf in monosyllable
		"נָעֳמִי":    "noʿŏmî",      // qamatz before hataf-qamatz
	})
}

func TestFurtivePatah(t *testing.T) {
	runCases(t, "furtive", map[string]string{
		"נֹחַ":       "nōaḥ",   // furtive on ḥet
		"רָקִיעַ":    "rāqîaʿ", // furtive on ayin after mater
		"גָּבֹהַּ":   "gābōah", // furtive on he-with-mappiq
	})
}

func TestDageshForte(t *testing.T) {
	runCases(t, "dagesh forte", map[string]string{
		"הִנֵּה":       "hinnēh",  // doubled nun (tsere+he stays separate)
		"מַגָּל":       "maggāl",  // doubled gimel
		"מִנְּזָר":     "minnəzār", // doubled nun, vocal shewa after
		"שַׁבָּת":      "šabbāt",   // doubled bet mid-word
	})
}

func TestDageshLene(t *testing.T) {
	// Word-initial dagesh (qal/lene) doesn't double.
	runCases(t, "dagesh lene", map[string]string{
		"בָּם":   "bām",  // initial bet+dagesh
		"דָּם":   "dām",  // initial dalet+dagesh
		"תָּם":   "tām",  // initial tav+dagesh
	})
}

func TestShewa(t *testing.T) {
	runCases(t, "shewa", map[string]string{
		"שְׁמֹר":     "šəmōr",    // vocal initial shewa
		"סַלְכָה":    "salkâ",    // silent shewa after patah, qamatz-he ligature
		"קָטַלְתְּ":  "qāṭalt",   // two final shewas both silent
	})
}

func TestDivineName(t *testing.T) {
	runCases(t, "divine name", map[string]string{
		"יְהוָה":  "yhwh",
		"יהוה":    "yhwh",
	})
}

func TestPunctuation(t *testing.T) {
	runCases(t, "punctuation", map[string]string{
		"שָׁלוֹם׃":    "šālôm:",     // sof passuq
		"בֶּן־אָדָם":  "ben-ʾādām",  // maqaf
	})
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{"Hello", "\\v 1 In the beginning"} {
		if got := Transliterate(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("שלום") {
		t.Error("expected true for Hebrew")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin")
	}
}

func TestEmpty(t *testing.T) {
	if got := Transliterate(""); got != "" {
		t.Errorf("empty input: got %q", got)
	}
}

// TestSBLCorpus exercises additional cases drawn from charlesLoder/
// hebrew-transliteration's default SBL Academic test suite, ensuring
// the rule coverage matches the reference implementation on the cases
// that don't depend on accent-driven long-hiriq/qubuts macrons.
func TestSBLCorpus(t *testing.T) {
	runCases(t, "sbl corpus", map[string]string{
		// full alphabet
		"אבגדהוזחטיכךלמםנןסעפףצץקרשת": "ʾbgdhwzḥṭykklmmnnsʿppṣṣqršt",
		// non-Hebrew passthrough
		"v1. רַעַל":   "v1. raʿal",
		"v1.\n רַעַל": "v1.\n raʿal",

		// BGDKPT — default SBL Academic does not mark spirantization
		"גָּדַל": "gādal",
		"חָג":   "ḥāg",
		"סַד":   "sad",
		"כָּמָר": "kāmār",
		"לֵךְ":  "lēk",
		"פֹּה":  "pōh",
		"אֶלֶף": "ʾelep",
		"מַת":   "mat",
		"אָרַשׂ": "ʾāraś",

		// Dagesh tests beyond the basics
		"בֹּסֶר":   "bōser",
		"מַסְגֵּר": "masgēr",

		// Vowel features
		"מִנְחָה":   "minḥâ",
		"וַיֻּגַּד": "wayyuggad",

		// Mater edges
		"יַיִן":    "yayin",   // yod-yod where second carries its own vowel
		"רִיֵם":    "riyēm",   // yod with own vowel is not a mater
		"בִּיטוֹן": "bîṭôn",   // bgdkpt letter with mater
		"כָּךְ":   "kāk",      // final shewa silent

		// 3MS suffix appears naturally from the rules
		"דְּבָרָיו": "dəbārāyw",
	})
}
