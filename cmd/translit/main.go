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

	"github.com/digitalbiblesociety/transliterate/script"
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
		fmt.Fprint(os.Stdout, rootHelp())
		return
	}
	switch args[0] {
	case "usfm":
		fmt.Fprint(os.Stdout, usfmHelp)
	case "bibles":
		fmt.Fprint(os.Stdout, biblesHelp)
	case "text":
		fmt.Fprint(os.Stdout, textHelp())
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n", args[0])
		fmt.Fprint(os.Stderr, rootHelp())
		os.Exit(2)
	}
}

// rootHelp / textHelp pull from the mode registry, which is populated
// by the script package's init() — must be evaluated at call time, not
// as a top-level var.
func rootHelp() string {
	return `translit — transliterate text from 30+ writing systems to Latin.

Usage:
  translit [flags] [words...]       one-shot text mode (auto-detect)
  translit usfm   -in <d> -out <d>  walk a directory of .usfm files
  translit bibles -root <d>         walk a tree of <bible>/usfm/ dirs
  translit help [subcommand]

Scripts are identified by ISO 15924 codes (case-insensitive):
  ` + strings.Join(script.Names(), ", ") + `.

` + modesHelp() + `
Run 'translit help <subcommand>' for subcommand-specific flags.
`
}

func textHelp() string {
	return `translit (text mode) — transliterate words or stdin.

Usage:
  translit [flags] [words...]
  echo "<input>" | translit [flags]

Flags:
  -script <code>   Force a specific script by ISO 15924 code (case-insensitive).
                   One of: ` + strings.Join(script.Names(), ", ") + `.
  -mode <name>     Alternate transliteration mode. No-op if the detected
                   script doesn't have one with this name.

` + modesHelp() + `
Legacy aliases (deprecated, still functional):
  -tashkeel   = -mode tashkeel
  -notones    = -mode atonal
  -phonetic   = -mode phonetic

If words are passed as arguments, each is transliterated on its own
output line. With no arguments, lines are read from stdin.
`
}
