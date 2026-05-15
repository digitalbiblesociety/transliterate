// Command translit transliterates text from 19 writing systems to the
// Latin alphabet. It has three modes:
//
//	translit [flags] [words...]       one-shot text mode (auto-detect)
//	translit usfm   [flags]           walk a directory of .usfm files
//	translit bibles [flags]           walk a tree of <bible>/usfm/ dirs
//
// Run `translit help` for the full flag list of each subcommand, or
// `translit help <subcommand>` for one in particular.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		runText(nil)
		return
	}
	switch os.Args[1] {
	case "usfm":
		runUSFM(os.Args[2:])
	case "bibles":
		runBibles(os.Args[2:])
	case "help", "-h", "--help":
		printHelp(os.Args[2:])
	default:
		runText(os.Args[1:])
	}
}

func printHelp(args []string) {
	if len(args) == 0 {
		fmt.Fprint(os.Stdout, rootHelp)
		return
	}
	switch args[0] {
	case "usfm":
		fmt.Fprint(os.Stdout, usfmHelp)
	case "bibles":
		fmt.Fprint(os.Stdout, biblesHelp)
	case "text":
		fmt.Fprint(os.Stdout, textHelp)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n", args[0])
		fmt.Fprint(os.Stderr, rootHelp)
		os.Exit(2)
	}
}

var rootHelp = `translit — transliterate text from 30+ writing systems to Latin.

Usage:
  translit [flags] [words...]       one-shot text mode (auto-detect)
  translit usfm   -in <d> -out <d>  walk a directory of .usfm files
  translit bibles -root <d>         walk a tree of <bible>/usfm/ dirs
  translit help [subcommand]

Scripts are identified by ISO 15924 codes (case-insensitive):
  ` + strings.Join(engineNames(), ", ") + `.

Run 'translit help <subcommand>' for subcommand-specific flags.
`

var textHelp = `translit (text mode) — transliterate words or stdin.

Usage:
  translit [flags] [words...]
  echo "<input>" | translit [flags]

Flags:
  -script <code>   Force a specific script by ISO 15924 code (case-insensitive).
                   One of: ` + strings.Join(engineNames(), ", ") + `.
  -tashkeel        For Arab input, use the tashkeel-aware engine instead of
                   the default ANETAC dictionary lookup.
  -notones         For Hani (Mandarin) or Yueh (Cantonese) input, strip tone
                   marks (Mandarin diacritics) or tone digits (Jyutping 1-6).

If words are passed as arguments, each is transliterated on its own
output line. With no arguments, lines are read from stdin.
`
