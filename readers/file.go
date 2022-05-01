package readers

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadFromFile opens and reads from file at filepath. It gracefully
// handles errors.
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

// ReadTextualData reads header and content from file.
func ReadTextualData(file io.Reader, filename string) (header string, content string, err error) {
	data, err := io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("failed to read from file %s: %v", filename, err)
		return
	}

	text := string(data)

	chunks := strings.Split(text, "\n\n")
	header = chunks[0]

	for i := 1; i < len(chunks)-1; i++ {
		chunk := chunks[i]
		chunk = strings.ReplaceAll(chunk, "\n", " ")
		content += chunk

		if i != len(text)-1 {
			content += "\n"
		}
	}
	return
}
