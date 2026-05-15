package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/digitalbiblesociety/transliterate/script"
)

var usfmHelp = `translit usfm — walk a directory of .usfm files.

Usage:
  translit usfm -in <source-dir> -out <dest-dir> [flags]

Flags:
  -in <dir>        Input directory containing .usfm files (required).
  -out <dir>       Output directory (created if missing; required).
  -script <name>   Force a specific script (default: auto-detect across files).
  -mode <name>     Alternate transliteration mode (see 'translit help' for list).
  -jobs N          Files to process in parallel (default: NumCPU).

USFM markers, ASCII, and content outside the detected script's block
pass through unchanged.
`

func runUSFM(args []string) {
	fs := flag.NewFlagSet("translit usfm", flag.ExitOnError)
	in := fs.String("in", "", "input directory containing .usfm files")
	out := fs.String("out", "", "output directory (created if missing)")
	scriptName := fs.String("script", "", "force a specific script (ISO 15924; default: auto-detect)")
	mf := registerModeFlags(fs)
	jobs := fs.Int("jobs", runtime.NumCPU(), "files to process in parallel")
	fs.Usage = func() { fmt.Fprint(os.Stderr, usfmHelp) }
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	mf.validate()

	if *in == "" || *out == "" {
		fmt.Fprint(os.Stderr, usfmHelp)
		os.Exit(2)
	}

	files, err := listUSFM(*in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "list input:", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "no .usfm files in", *in)
		os.Exit(1)
	}

	var eng *script.Engine
	if *scriptName != "" {
		eng = script.ByName(*scriptName)
		if eng == nil {
			fmt.Fprintf(os.Stderr, "unknown script %q\nvalid: %s\n", *scriptName, strings.Join(script.Names(), ", "))
			os.Exit(2)
		}
	} else {
		eng = detectAcrossFiles(*in, files)
		if eng == nil {
			fmt.Fprintln(os.Stderr, "no recognized script detected in", *in)
			os.Exit(1)
		}
	}

	transFn := script.ResolveMode(eng, mf.effective())

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "create output dir:", err)
		os.Exit(1)
	}

	start := time.Now()
	var done, failed atomic.Int64
	sem := make(chan struct{}, *jobs)
	var wg sync.WaitGroup
	for _, name := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(name string) {
			defer wg.Done()
			defer func() { <-sem }()
			src := filepath.Join(*in, name)
			dst := filepath.Join(*out, name)
			if err := processUSFMFile(transFn, src, dst); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
				failed.Add(1)
				return
			}
			n := done.Add(1)
			fmt.Fprintf(os.Stderr, "[%d/%d] %s\n", n, len(files), name)
		}(name)
	}
	wg.Wait()

	fmt.Fprintf(os.Stderr, "[%s] %d ok, %d failed in %s\n",
		eng.Name, done.Load(), failed.Load(), time.Since(start).Round(time.Millisecond))
	if failed.Load() > 0 {
		os.Exit(1)
	}
}

func listUSFM(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".usfm") {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

// detectAcrossFiles reads up to 32 KB from each file and tallies block
// matches across every engine. Returns the engine with the highest
// total. Used by `usfm` and `bibles` subcommands instead of detecting
// per-file because a directory of USFM files all describe the same
// translation and should be processed uniformly.
func detectAcrossFiles(dir string, files []string) *script.Engine {
	engines := script.Engines()
	counts := make([]int, len(engines))
	for _, name := range files {
		f, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		buf := make([]byte, 32<<10)
		n, _ := io.ReadFull(f, buf)
		f.Close()
		for _, r := range string(buf[:n]) {
			for i := range engines {
				if engines[i].Contains(r) {
					counts[i]++
					break
				}
			}
		}
	}
	best, bestN := -1, 0
	for i, c := range counts {
		if c > bestN {
			best, bestN = i, c
		}
	}
	if best < 0 {
		return nil
	}
	return &engines[best]
}

// processUSFMFile is a thin wrapper around the public
// script.TransliterateFile helper; kept as a local alias so the
// existing call sites don't churn.
func processUSFMFile(transliterate func(string) string, src, dst string) error {
	return script.TransliterateFile(transliterate, src, dst)
}
