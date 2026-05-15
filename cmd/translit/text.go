package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/digitalbiblesociety/transliterate/script"
)

func runText(args []string) {
	fs := flag.NewFlagSet("translit", flag.ExitOnError)
	scriptName := fs.String("script", "", "force a specific script (ISO 15924; default: auto-detect)")
	mf := registerModeFlags(fs)
	fs.Usage = func() { fmt.Fprint(os.Stderr, textHelp()) }
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	mf.validate()

	var forced *script.Engine
	if *scriptName != "" {
		forced = script.ByName(*scriptName)
		if forced == nil {
			fmt.Fprintf(os.Stderr, "unknown script %q\nvalid: %s\n", *scriptName, strings.Join(script.Names(), ", "))
			os.Exit(2)
		}
	}

	mode := mf.effective()
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
			eng = script.Detect(in)
		}
		if eng == nil {
			fmt.Fprintln(w, in)
			return
		}
		fmt.Fprintln(w, script.ResolveMode(eng, mode)(in))
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
