package readers

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// AskForConfirmation prints message to stdout and presents user
// with the boolean choice. It returns a bool indicating whether the user
// agreed or rejected or and error, when user's response is invalid.
func AskForConfirmation(message string, defaultYes bool) (bool, error) {
	yesAnswers := make(map[string]bool)
	yesAnswers["Y\n"] = true
	yesAnswers["y\n"] = true

	noAnswers := make(map[string]bool)
	noAnswers["N\n"] = true
	noAnswers["n\n"] = true

	if defaultYes {
		yesAnswers["\n"] = true
		fmt.Printf(message + " [Y/n] ")
	} else {
		noAnswers["\n"] = true
		fmt.Printf(message + " [y/N] ")
	}

	reader := bufio.NewReader(os.Stdin)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false, errors.WithStack(err)
	}

	if yesAnswers[response] {
		return true, nil
	} else if noAnswers[response] {
		return false, nil
	}

	return false, errors.New("unknown option passed")
}
