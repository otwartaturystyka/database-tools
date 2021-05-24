package readers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// AskForConfirmation prints message to stdout and prompts the user
// to say yes or no. It returns a bool indicating whether the user
// agreed or rejected or and error, when user's response is invalid.
func AskForConfirmation(r io.Reader, w io.Writer, message string, defaultYes bool) (bool, error) {
	yesAnswers := make(map[string]bool)
	yesAnswers["Y\n"] = true
	yesAnswers["y\n"] = true

	noAnswers := make(map[string]bool)
	noAnswers["N\n"] = true
	noAnswers["n\n"] = true

	if defaultYes {
		yesAnswers["\n"] = true
		fmt.Fprint(w, message+" [Y/n] ")
	} else {
		noAnswers["\n"] = true
		fmt.Fprint(w, message+" [y/N] ")
	}

	reader := bufio.NewReader(r)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read from reader: %w", err)
	}

	if yesAnswers[response] {
		return true, nil
	} else if noAnswers[response] {
		return false, nil
	}

	return false, errors.New("unknown option passed")
}
