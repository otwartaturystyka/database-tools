package readers

import (
	"bytes"
	"os"
	"testing"
)

func TestAskForConfirmation(t *testing.T) {
	t.Run("valid agreement", func(t *testing.T) {
		var buffer bytes.Buffer
		buffer.WriteString("y\n")

		want := true

		ignoredOut, _ := os.Open(os.DevNull)
		got, err := AskForConfirmation(&buffer, ignoredOut, "question?", false)
		if err != nil {
			t.Errorf("\nfailed to get response: %v\n", err)
		}

		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})

	t.Run("valid disagreement", func(t *testing.T) {
		buffer := bytes.Buffer{}
		buffer.WriteString("n\n")

		want := false

		ignoredOut, _ := os.Open(os.DevNull)
		got, err := AskForConfirmation(&buffer, ignoredOut, "question?", false)
		if err != nil {
			t.Errorf("\nfailed to get response: %v\n", err)
		}

		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})

	t.Run("invalid text", func(t *testing.T) {
		var buffer bytes.Buffer
		buffer.WriteString("some invalid text\n")

		ignoredOut, _ := os.Open(os.DevNull)
		got, err := AskForConfirmation(&buffer, ignoredOut, "question?", false)
		if err == nil {
			t.Errorf("wanted error, got %t", got)
		}
	})
}
