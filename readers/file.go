package readers

import (
	"fmt"
	"io"
	"os"
)

// ReadFromFile opens and reads from file at filepath. It gracefully handles
// errors.
func ReadFromFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", filepath, err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read all from file %s: %w", filepath, err)
	}

	return fileContent, nil
}
