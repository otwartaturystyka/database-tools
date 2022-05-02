// Package formatters provides utilities for formatting text read from source
// files into a nicely looking text that lands in the app.
package formatters

import "strings"

// ToContent converts whitespace and newline formatting so that it looks good
// in the mobile app.
//
// Used for:
// - Quick info
// - Overview
// In sections, places and trails.
func ToContent(text string) (formattedText string) {
	chunks := strings.Split(text, "\n\n")

	for i, chunk := range chunks {
		chunk = strings.ReplaceAll(chunk, "\n", " ")

		if i != len(chunks)-1 {
			chunk += "\n\n"
		} else {
			chunk = strings.TrimSuffix(chunk, " ")
		}

		formattedText += chunk
	}

	return formattedText
}

// ToSection reads a section (consisting of header and content) from r.
// - Quick info
// - Overview
func ToSection(text string) (header string, content string) {
	chunks := strings.Split(text, "\n\n")
	header = chunks[0]

	for i := 1; i < len(chunks); i++ {
		chunk := chunks[i]
		chunk = strings.ReplaceAll(chunk, "\n", " ")

		if i != len(chunks)-1 {
			chunk += "\n\n"
		} else {
			chunk = strings.TrimSuffix(chunk, " ")
		}

		content += chunk
	}

	return
}
