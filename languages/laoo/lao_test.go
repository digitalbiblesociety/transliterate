package laoo

import "testing"

func TestPlaceNames(t *testing.T) {
	cases := map[string]string{
		"ລາວ":       "lāo",
		"ວຽງຈັນ":    "viangchan",
		"ຫຼວງພະບາງ": "luangphabāng",
		"ນະຄອນຫຼວງ": "nakhǭnluang",
		"ພາສາລາວ":   "phāsālāo",
		"ໂຮງຮຽນ":    "hōnghian",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPrepositionedVowels(t *testing.T) {
	cases := map[string]string{
		"ເດັກ":  "dek",
		"ແມວ":   "mǣo",
		"ໂຮງ":   "hōng",
		"ໄນ້":   "nai",
		"ໃຫຍ່":  "nyai",
		"ເວົ້າ": "vao",
		"ເຫຼົາ": "lao",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSilentHoSung(t *testing.T) {
	cases := map[string]string{
		"ຫຼວງ": "luang",
		"ຫງາຍ": "ngāy",
		"ຫມາ":  "mā",
		"ຫນ້າ": "nā",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFinalRemapping(t *testing.T) {
	cases := map[string]string{
		"ກາດ": "kāt",
		"ກາບ": "kāp",
		"ກາວ": "kāo",
		"ກາງ": "kāng",
		"ການ": "kān",
		"ກາມ": "kām",
		"ກາກ": "kāk",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestVowelSigns(t *testing.T) {
	cases := map[string]string{
		"ກະ": "ka",
		"ກາ": "kā",
		"ກິ": "ki",
		"ກີ": "kī",
		"ກຸ": "ku",
		"ກູ": "kū",
		"ກຶ": "kư",
		"ກື": "kư̄",
		"ກຳ": "kam",
		"ກໍ": "kǭ",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestToneMarksDropped(t *testing.T) {
	cases := map[string]string{
		"ປ່າ": "pā",
		"ປ້າ": "pā",
		"ປ໊າ": "pā",
		"ປ໋າ": "pā",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPaliStackCancellation(t *testing.T) {
	if got := Transliterate("ສິລ໌"); got != "si" {
		t.Errorf("Transliterate(ສິລ໌) = %q, want %q", got, "si")
	}
}

func TestDigits(t *testing.T) {
	cases := map[string]string{
		"໑໒໓໔໕":   "12345",
		"໐໖໗໘໙":   "06789",
		"ປີ ໒໐໒໕": "pī 2025",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{"Hello", "1:1", "", "Lao: ລາວ"} {
		got := Transliterate(s)
		if !Contains(s) && got != s {
			t.Errorf("non-Lao %q: got %q, want passthrough", s, got)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("ລາວ") {
		t.Error("expected true for Lao string")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin string")
	}
	if Contains("") {
		t.Error("expected false for empty string")
	}
}
