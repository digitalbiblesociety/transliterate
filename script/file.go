package script

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// TransliterateFile reads src, applies the transliterate function to
// each line's body (without the trailing newline), and writes the
// result to dst. Writes go through a `.tmp` sibling that's renamed on
// success — partially-written outputs are never left behind under dst.
//
// Designed for USFM files: marker lines (`\v 1`, `\c 1`, etc.) and
// pure-ASCII content pass through unchanged because the per-script
// Transliterate functions only act on their own block. Streamed via
// bufio so a 64-KiB read buffer copes with the largest scripture
// books without loading them whole.
func TransliterateFile(transliterate func(string) string, src, dst string) error {
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
