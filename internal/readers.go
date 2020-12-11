package internal

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// ReadFromFile opens and reads from file at filepath. It gracefully
// handles errors.
func readFromFile(filepath string) ([]byte, error) {
	nameFile, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Errorf("failed to open file %s", filepath)
	}
	defer nameFile.Close()

	name, err := ioutil.ReadAll(nameFile)
	if err != nil {
		return nil, errors.Errorf("failed to read contents from file %s", filepath)
	}

	return name, nil
}

// ReadTextualData reads header and content from file.
func readTextualData(file *os.File) (header string, content string, err error) {
	reader := bufio.NewReader(file)
	header, err = reader.ReadString('\n')
	if err != nil {
		err = errors.Errorf("failed to read header (line 1)")
		return
	}

	reader.ReadString('\n')
	if err != nil {
		err = errors.Errorf("failed to read 3-slash divider (line 2)")
		return
	}

	content, err = reader.ReadString('\n')
	if err != nil {
		err = errors.Errorf("failed to read content (line 3)")
		return
	}

	return
}
