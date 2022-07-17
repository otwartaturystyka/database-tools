package readers

import (
	"fmt"
	"io"
	"os"
)

// ReadLocalized returns a map of languages to texts .
func ReadLocalized(filepath string) (map[string]string, error) {
	entries, err := os.ReadDir("./")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}

	return nil, nil
}

// ReadFromFile opens and reads from file at filepath. It gracefully handles
// errors.
func ReadFromFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return fileContent, nil
}
