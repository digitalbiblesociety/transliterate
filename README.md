![transliterate](./assets/detector.svg)

![transliterate](./assets/transliterate-logo.svg)

Pure-Go transliteration of 30+ writing systems to the Latin alphabet, built for Bible-translation workflows utilizing the [USFM format](https://ubsicap.github.io/usfm/about/index.html) and the [Aeneas Project](https://github.com/readbeyond/aeneas) like [Scripture App Builder](https://software.sil.org/scriptureappbuilder/). But usable for any text-processing pipeline.

To request an additional language, offer transliteration corrections/preferences, please create an [issue](https://github.com/digitalbiblesociety/transliterate/issues).

## Supported scripts

Packages are named by their [ISO 15924](https://www.unicode.org/iso15924/iso15924-codes.html)
four-letter code _lowercased for go_. All live under `github.com/digitalbiblesociety/transliterate/languages/<code>`.

| Pkg    | Script                | Languages (sample)                                                    | Scheme               |
| ------ | --------------------- | --------------------------------------------------------------------- | -------------------- |
| `arab` | Arabic                | Arabic, Persian, Urdu (script)                                        | ANETAC + rules       |
| `armn` | Armenian              | Armenian                                                              | ISO 9985             |
| `bali` | Balinese              | Balinese, Kawi                                                        | ISO 15919            |
| `batk` | Batak                 | Toba, Karo, Mandailing, Pakpak, Simalungun                            | ISO 15919            |
| `beng` | Bengali               | Bengali, Assamese                                                     | ISO 15919            |
| `brah` | Brahmi                | Aśokan Prakrits, early Buddhist/Jain (historical)                     | ISO 15919            |
| `bugi` | Buginese / Lontara    | Bugis, Makassar                                                       | ISO 15919            |
| `cans` | Canadian Syllabics    | Cree, Inuktitut, Ojibwe, Naskapi, Blackfoot                           | Unicode-name derived |
| `cher` | Cherokee              | Cherokee                                                              | Sequoyah             |
| `cyrl` | Cyrillic              | Russian, Ukrainian, Bulgarian, Serbian, Mongolian                     | ISO 9 (1995)         |
| `deva` | Devanagari            | Hindi, Sanskrit, Marathi, Nepali, Bhojpuri                            | ISO 15919            |
| `ethi` | Ethiopic / Ge'ez      | Amharic, Tigrinya, Ge'ez liturgical, Tigre                            | BGN/PCGN             |
| `geor` | Georgian              | Georgian (Mkhedruli, Asomtavruli)                                     | BGN/PCGN             |
| `grek` | Greek                 | Ancient/Koine/Modern Greek (polytonic)                                | SBL                  |
| `gujr` | Gujarati              | Gujarati                                                              | ISO 15919            |
| `guru` | Gurmukhi              | Punjabi                                                               | ISO 15919            |
| `hang` | Hangul                | Korean                                                                | Revised Romanization |
| `hani` | Han                   | Mandarin (CJK U+4E00..U+9FFF) — phrase-aware polyphones               | Hanyu Pinyin (tonal) |
| `hebr` | Hebrew                | Hebrew (pointed and unpointed)                                        | SBL                  |
| `java` | Javanese              | Javanese, Kawi (Old Javanese)                                         | ISO 15919            |
| `jpan` | Japanese (Hira+Kata+Kanji) | Japanese — kana + kanji via kagome morphological analyzer        | Hepburn              |
| `khmr` | Khmer                 | Khmer                                                                 | Simplified           |
| `knda` | Kannada               | Kannada, Tulu, Konkani                                                | ISO 15919            |
| `lana` | Tai Tham / Lanna      | Northern Thai, Tai Lue, Khün, Lao Tham                                | ISO 15919            |
| `laoo` | Lao                   | Lao                                                                   | BGN/PCGN             |
| `mlym` | Malayalam             | Malayalam                                                             | ISO 15919            |
| `modi` | Modi                  | Marathi (historical)                                                  | ISO 15919            |
| `mymr` | Myanmar / Burmese     | Burmese, Shan, Mon                                                    | BGN/PCGN             |
| `newa` | Newa / Prachalit      | Nepal Bhasa (Newari)                                                  | ISO 15919            |
| `orya` | Oriya                 | Odia                                                                  | ISO 15919            |
| `shrd` | Sharada               | Kashmiri Śaiva, Sanskrit (historical)                                 | ISO 15919            |
| `sinh` | Sinhala               | Sinhala                                                               | ISO 15919            |
| `sund` | Sundanese             | Sundanese                                                             | ISO 15919            |
| `syrc` | Syriac                | Classical Syriac (Peshitta, Eastern Christian)                        | ISO 233-3            |
| `taml` | Tamil                 | Tamil                                                                 | ISO 15919            |
| `telu` | Telugu                | Telugu                                                                | ISO 15919            |
| `thai` | Thai                  | Thai                                                                  | RTGS (PyThaiNLP port)|
| `tibt` | Tibetan               | Tibetan                                                               | Wylie + THL Phonetic |
| `tirh` | Tirhuta               | Maithili                                                              | ISO 15919            |
| `yueh` | Han (Cantonese)       | Cantonese / Yue Chinese                                               | Jyutping (tone digit)|

## Install

Library:

```sh
go get github.com/digitalbiblesociety/transliterate
```

CLI:

```sh
go install github.com/digitalbiblesociety/transliterate/cmd/translit@latest
```

## CLI

A single `translit` binary covers every workflow. Auto-detects the script of
the input; `-script <code>` (ISO 15924, case-insensitive) forces a specific one.

```sh
# One-shot (auto-detect):
$ translit اميس باراندياران العراق
Ames
Barandiaran
Alaraq

$ translit "Ἰησοῦς Χριστός"
Iēsous Christos

$ translit "ಆದಿಯಲ್ಲಿ ದೇವರು"
ādiyalli dēvaru

# Read from stdin:
$ echo "Москва" | translit
Moskva

# Force a specific script by ISO 15924 code:
$ translit -script Hebr "בְּרֵאשִׁית"
bərēʾšiyt

# For vocalized Arabic, use the tashkeel-aware engine:
$ translit -mode tashkeel "يَعْقُوبَ"
yaequb

# Chinese auto-detects as Mandarin (tonal); read as Cantonese explicitly:
$ translit "你好"
nǐ hǎo
$ translit -mode atonal "你好"
ni hao
$ translit -script Yueh "你好"
nei5 hou2

# Tibetan defaults to Wylie; -mode phonetic switches to THL Simplified
# Phonetic (approximate Lhasa pronunciation, ideal for forced alignment):
$ translit "བཀྲ་ཤིས་བདེ་ལེགས"
bkra shis bde legs
$ translit -mode phonetic "བཀྲ་ཤིས་བདེ་ལེགས"
tra shi de lek

# USFM directory walker (auto-detects per directory):
$ translit usfm -in ./source-usfm -out ./latin-usfm -jobs 8

# Multi-Bible directory walker (auto-detects per Bible):
$ translit bibles -root /path/to/bibles -force
```

Run `translit help` (or `translit help <subcommand>`) for the full flag list.

`make build` produces the binary at `./bin/translit`.

## Datasets and attribution

Third-party data and algorithms in this repo:
- **ANETAC** (MIT) — Arabic named-entity dictionary.
- **Unicode Unihan Database** (Unicode Terms of Use) — Cantonese
  Jyutping readings (Mandarin moved to go-pinyin, see below).
- **Unicode character names** (Unicode Terms of Use) — Cansyl and
  Cherokee tables.
- **PyThaiNLP** (Apache 2.0) — algorithmic basis for the Thai RTGS
  port; no source code is bundled.
- **Aksharamukha** (MIT) — codepoint romanization tables for Brahmi,
  Sharada, Modi, Tirhuta, Newa, and Tai Tham.
- **kagome** + **kagome-dict** (MIT) — Japanese morphological analyzer
  and IPA dictionary. Vendored under `internal/kagome/` rather than
  pulled via `go.mod` so the build has no external Go module
  dependencies; bundled `ipa.dict` is from mecab-ipadic-2.7.0-20070801
  (NAIST/ICOT terms).
- **go-pinyin** + **phrase-pinyin-data** (both MIT) — Mandarin per-
  character readings and ~47k phrase dictionary for polyphone
  disambiguation (中 reads zhōng in 中国 but zhòng in 击中). Vendored
  under `internal/pinyin/`.

Full attribution in [NOTICE](NOTICE).

All other script tables are based on public international standards
(ISO, BGN/PCGN, RR, RTGS, Hepburn, SBL, Wylie).

> Note on vendoring: we keep large third-party Go projects in-tree under
> `internal/` (rather than as `go.mod` dependencies) to keep the build
> hermetic and to reduce the supply-chain surface — `go build` should not
> need to pull anything from the module proxy at compile time.

## License

Freely available under the [MIT License](LICENSE)

<p align="center">
  <img src="https://raw.githubusercontent.com/digitalbiblesociety/transliterate/refs/heads/main/assets/dbs-logo-half.svg" />
</p>