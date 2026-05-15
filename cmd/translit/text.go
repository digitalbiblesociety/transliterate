package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func runText(args []string) {
	fs := flag.NewFlagSet("translit", flag.ExitOnError)
	scriptName := fs.String("script", "", "force a specific script (ISO 15924; default: auto-detect)")
	tashkeel := fs.Bool("tashkeel", false, "for Arab input, use the tashkeel-aware engine (no-op for other scripts)")
	notones := fs.Bool("notones", false, "for Hani/Yueh input, strip tone marks (Mandarin) or tone digits (Cantonese)")
	fs.Usage = func() { fmt.Fprint(os.Stderr, textHelp) }
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	var forced *engine
	if *scriptName != "" {
		forced = engineByName(*scriptName)
		if forced == nil {
			fmt.Fprintf(os.Stderr, "unknown script %q\nvalid: %s\n", *scriptName, strings.Join(engineNames(), ", "))
			os.Exit(2)
		}
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	emit := func(in string) {
		in = strings.TrimRight(in, "\r\n")
		if in == "" {
			fmt.Fprintln(w)
			return
		}
		eng := forced
		if eng == nil {
			eng = detect(in)
		}
		if eng == nil {
			fmt.Fprintln(w, in)
			return
		}
		fn := eng.transliterate
		if *tashkeel && eng.tashkeel != nil {
			fn = eng.tashkeel
		}
		if *notones && eng.atonal != nil {
			fn = eng.atonal
		}
		fmt.Fprintln(w, fn(in))
	}

	if fs.NArg() > 0 {
		for _, a := range fs.Args() {
			emit(a)
		}
		return
	}

	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if line != "" {
			emit(line)
		}
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "read stdin:", err)
			os.Exit(1)
		}
	}
}
