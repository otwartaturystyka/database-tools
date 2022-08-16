package readers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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

// ReadLocalizedFiles reads contents of filename in all available languages.
//
// It lists names of directories in cwd, and then reads contents of filename in
// every of these directories.
func ReadLocalizedFiles(filename string) (map[string]string, error) {
	dirs, err := os.ReadDir("content")
	if err != nil {
		return nil, fmt.Errorf("read localized files: %v", err)
	}

	contents := make(map[string]string)
	for _, dir := range dirs {
		lang := dir.Name()
		filepath := filepath.Join("content", lang, filename)
		content, err := ReadFromFile(filepath)
		if err != nil {
			fmt.Printf("error reading file: %v\n", err)
			continue // it's okay
		}

		contents[lang] = string(content)
	}

	_, isEnglish := contents["en"]
	_, isPolish := contents["pl"]

	if !isEnglish && !isPolish {
		fmt.Println("no english nor polish for", filename)
		return nil, err
	}

	if _, ok := contents["pl"]; !ok {
		return nil, fmt.Errorf("no pl translation for %s", filename)
	}

	return contents, nil
}
