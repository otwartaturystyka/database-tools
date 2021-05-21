package readers

import (
	"bytes"
	"testing"
)

func TestReadTextualDataNormal(t *testing.T) {
	var buffer bytes.Buffer
	wantHeader := "Example header"
	wantContent := "Example content"

	buffer.WriteString(wantHeader + "\n---\n" + wantContent + "\n")

	gotHeader, gotContent, err := ReadTextualData(&buffer, "testfile")
	if err != nil {
		t.Error(err)
	}
	if gotHeader != wantHeader {
		t.Errorf("got %q, want %q", gotHeader, wantHeader)
	}

	if gotContent != wantContent {
		t.Errorf("got %q, want %q", gotContent, wantContent)
	}
}

func TestReadTextualDataTooLong(t *testing.T) {
	var buffer bytes.Buffer
	wantHeader := "Example header"
	wantContent := "Example content"

	buffer.WriteString(wantHeader + "\n---\n" + wantContent + "\nunused\nlines\n")

	gotHeader, gotContent, err := ReadTextualData(&buffer, "testfile")
	if err != nil {
		t.Error(err)
	}
	if gotHeader != wantHeader {
		t.Errorf("got %q, want %q", gotHeader, wantHeader)
	}

	if gotContent != wantContent {
		t.Errorf("got %q, want %q", gotContent, wantContent)
	}
}
