package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/digitalbiblesociety/transliterate/script"
)


var biblesHelp = `translit bibles — walk a tree of <bible>/usfm/ directories.

Usage:
  translit bibles -root <dir> [flags]

For each <bible>/usfm/ inside -root, produces a sibling <bible>/<out-name>/
with every .usfm file transliterated. Auto-detects the script per Bible.

Flags:
  -root <dir>      Root containing <bible>/usfm/ subdirs (required).
  -force           Overwrite existing output dirs (default: skip).
  -out-name <name> Sibling output dir name (default: usfm-transliterate).
  -mode <name>     Alternate transliteration mode (see 'translit help' for list).
  -jobs N          Bibles to process in parallel (default: NumCPU).
`

type bibleResult struct {
	bible   string
	script  string
	files   int
	errs    int
	dur     time.Duration
	skipped bool
}

func runBibles(args []string) {
	fs := flag.NewFlagSet("translit bibles", flag.ExitOnError)
	root := fs.String("root", "", "root containing <bible>/usfm/ subdirs")
	force := fs.Bool("force", false, "overwrite existing output dirs")
	outName := fs.String("out-name", "usfm-transliterate", "sibling output dir name")
	mf := registerModeFlags(fs)
	jobs := fs.Int("jobs", runtime.NumCPU(), "Bibles to process in parallel")
	fs.Usage = func() { fmt.Fprint(os.Stderr, biblesHelp) }
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	mf.validate()

	if *root == "" {
		fmt.Fprint(os.Stderr, biblesHelp)
		os.Exit(2)
	}

	bibles, err := findBibles(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scan:", err)
		os.Exit(1)
	}

	results := make(chan bibleResult, len(bibles))
	sem := make(chan struct{}, *jobs)
	var wg sync.WaitGroup
	var processed, skipped, failed atomic.Int64

	for _, bible := range bibles {
		wg.Add(1)
		sem <- struct{}{}
		go func(bible string) {
			defer wg.Done()
			defer func() { <-sem }()

			inDir := filepath.Join(*root, bible, "usfm")
			outDir := filepath.Join(*root, bible, *outName)

			if !*force {
				if _, err := os.Stat(outDir); err == nil {
					skipped.Add(1)
					results <- bibleResult{bible: bible, skipped: true}
					return
				}
			}

			r := processBible(inDir, outDir, mf.effective())
			r.bible = bible
			if r.script == "" {
				skipped.Add(1)
			} else if r.errs > 0 {
				failed.Add(1)
			} else {
				processed.Add(1)
			}
			results <- r
		}(bible)
	}

	go func() { wg.Wait(); close(results) }()

	start := time.Now()
	for r := range results {
		switch {
		case r.skipped && r.script == "":
			fmt.Printf("- skip  %-12s already done\n", r.bible)
		case r.script == "":
			fmt.Printf("- skip  %-12s no script detected\n", r.bible)
		case r.errs > 0:
			fmt.Printf("x fail  %-12s [%s] %d files, %d errors in %s\n",
				r.bible, r.script, r.files, r.errs, r.dur.Round(time.Millisecond))
		default:
			fmt.Printf("> done  %-12s [%-10s] %d files in %s\n",
				r.bible, r.script, r.files, r.dur.Round(time.Millisecond))
		}
	}
	fmt.Printf("\nsummary: %d processed, %d skipped, %d failed in %s\n",
		processed.Load(), skipped.Load(), failed.Load(), time.Since(start).Round(time.Millisecond))
}

func findBibles(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		usfm := filepath.Join(root, e.Name(), "usfm")
		fi, err := os.Stat(usfm)
		if err != nil || !fi.IsDir() {
			continue
		}
		ents, err := os.ReadDir(usfm)
		if err != nil {
			continue
		}
		for _, f := range ents {
			if strings.HasSuffix(strings.ToLower(f.Name()), ".usfm") {
				out = append(out, e.Name())
				break
			}
		}
	}
	return out, nil
}

func processBible(inDir, outDir, mode string) bibleResult {
	start := time.Now()
	var r bibleResult

	files, err := listUSFM(inDir)
	if err != nil || len(files) == 0 {
		return r
	}

	eng := detectAcrossFiles(inDir, files)
	if eng == nil {
		return r
	}
	r.script = eng.Name

	transFn := script.ResolveMode(eng, mode)

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		r.errs++
		r.dur = time.Since(start)
		return r
	}

	sem := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup
	var errs atomic.Int64
	for _, name := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(name string) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := processUSFMFile(transFn, filepath.Join(inDir, name), filepath.Join(outDir, name)); err != nil {
				errs.Add(1)
			}
		}(name)
	}
	wg.Wait()

	r.files = len(files)
	r.errs = int(errs.Load())
	r.dur = time.Since(start)
	return r
}
