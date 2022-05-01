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

// ReadSection reads a section (consisting of header and content) from file.
func ReadSection(file io.Reader, filename string) (header string, content string, err error) {
	data, err := io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("failed to read from file %s: %v", filename, err)
		return
	}

	chunks := strings.SplitN(string(data), "\n\n", 2)
	header = chunks[0]
	content = strings.TrimSuffix(chunks[1], "\n")

	return
}
