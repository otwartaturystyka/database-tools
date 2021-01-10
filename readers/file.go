package readers

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// ReadFromFile opens and reads from file at filepath. It gracefully
// handles errors.
func ReadFromFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", filepath)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Errorf("failed to read contents from file %s", filepath)
	}

	return fileContent, nil
}

// ReadTextualData reads header and content from file.
func ReadTextualData(file *os.File) (header string, content string, err error) {
	reader := bufio.NewReader(file)
	header, err = reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		} else {
			err = errors.Errorf("failed to read header (line 1) from file %s: %v", file.Name(), err)
			return
		}
	}
	header = strings.TrimSuffix(header, "\n")

	_, err = reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		} else {
			err = errors.Errorf("failed to read 3-slash divider (line 2) from file %s: %v", file.Name(), err)
			return
		}
	}

	content, err = reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		} else {
			err = errors.Errorf("failed to read content (line 3) from file %s: %v", file.Name(), err)
			return
		}
	}

	return
}