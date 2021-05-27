package readers

import (
	"bufio"
	"errors"
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
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents from file %s: %w", filepath, err)
	}

	return fileContent, nil
}

// ReadTextualData reads header and content from file.
func ReadTextualData(file io.Reader, filename string) (header string, content string, err error) {
	reader := bufio.NewReader(file)
	header, err = reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		} else {
			err = fmt.Errorf("failed to read header (line 1) from file %s: %w", filename, err)
			return
		}
	}
	header = strings.TrimSuffix(header, "\n")

	_, err = reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		} else {
			err = fmt.Errorf("failed to read 3-slash divider (line 2) from file %s: %w", filename, err)
			return
		}
	}

	content, err = reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		} else {
			err = fmt.Errorf("failed to read content (line 3) from file %s: %w", filename, err)
			return
		}
	}
	content = strings.TrimSuffix(content, "\n")

	return
}
