package readers

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadFromFile opens and reads from file at filepath. It gracefully handles
// errors.
func ReadFromFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %v", filepath, err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read contents from file %s: %v", filepath, err)
	}

	return fileContent, nil
}

// ReadSection reads a section (consisting of header and content) from r.
func ReadSection(r io.Reader) (header string, content string, err error) {
	data, err := io.ReadAll(r)
	if err != nil {
		err = fmt.Errorf("failed to read from file: %v", err)
		return
	}

	text := string(data)

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
