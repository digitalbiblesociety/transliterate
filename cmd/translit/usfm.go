package main

import (
	"bufio"
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
)

var usfmHelp = `translit usfm — walk a directory of .usfm files.

Usage:
  translit usfm -in <source-dir> -out <dest-dir> [flags]

Flags:
  -in <dir>        Input directory containing .usfm files (required).
  -out <dir>       Output directory (created if missing; required).
  -script <name>   Force a specific script (default: auto-detect across files).
  -jobs N          Files to process in parallel (default: NumCPU).

USFM markers, ASCII, and content outside the detected script's block
pass through unchanged.
`

func runUSFM(args []string) {
	fs := flag.NewFlagSet("translit usfm", flag.ExitOnError)
	in := fs.String("in", "", "input directory containing .usfm files")
	out := fs.String("out", "", "output directory (created if missing)")
	scriptName := fs.String("script", "", "force a specific script (ISO 15924; default: auto-detect)")
	tashkeel := fs.Bool("tashkeel", false, "for Arab input, use the tashkeel-aware engine")
	notones := fs.Bool("notones", false, "for Hani/Yueh input, strip tone marks/digits")
	jobs := fs.Int("jobs", runtime.NumCPU(), "files to process in parallel")
	fs.Usage = func() { fmt.Fprint(os.Stderr, usfmHelp) }
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

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

	var eng *engine
	if *scriptName != "" {
		eng = engineByName(*scriptName)
		if eng == nil {
			fmt.Fprintf(os.Stderr, "unknown script %q\nvalid: %s\n", *scriptName, strings.Join(engineNames(), ", "))
			os.Exit(2)
		}
	} else {
		eng = detectAcrossFiles(*in, files)
		if eng == nil {
			fmt.Fprintln(os.Stderr, "no recognized script detected in", *in)
			os.Exit(1)
		}
	}

	transFn := eng.transliterate
	if *tashkeel && eng.tashkeel != nil {
		transFn = eng.tashkeel
	}
	if *notones && eng.atonal != nil {
		transFn = eng.atonal
	}

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
		eng.name, done.Load(), failed.Load(), time.Since(start).Round(time.Millisecond))
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
// matches across every engine. Returns the engine with the highest total.
func detectAcrossFiles(dir string, files []string) *engine {
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
				if engines[i].inBlock(r) {
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

func processUSFMFile(transliterate func(string) string, src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	tmp := dst + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	bw := bufio.NewWriter(out)
	r := bufio.NewReaderSize(in, 1<<16)
	for {
		line, err := r.ReadString('\n')
		if line != "" {
			hasNL := strings.HasSuffix(line, "\n")
			body := line
			if hasNL {
				body = line[:len(line)-1]
			}
			bw.WriteString(transliterate(body))
			if hasNL {
				bw.WriteByte('\n')
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			out.Close()
			os.Remove(tmp)
			return fmt.Errorf("read %s: %w", src, err)
		}
	}
	if err := bw.Flush(); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dst)
}
