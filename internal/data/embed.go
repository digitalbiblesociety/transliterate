// Package data embeds the ANETAC parallel corpus.
//
// Source: https://github.com/HadjAmeur/ANETAC-Dataset
// "Arabic Named Entity Transliteration and Classification" — 79,924 Arabic
// named entities paired line-for-line with their English transliterations.
package data

import "embed"

//go:embed anetac/*.ar anetac/*.en
var Files embed.FS

// Splits lists the three corpus splits available in the embedded FS.
var Splits = []string{"train", "dev", "test"}
