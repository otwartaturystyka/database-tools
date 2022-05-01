package readers

import (
	"bytes"
	"testing"
)

func TestReadSectionNormal(t *testing.T) {
	var buffer bytes.Buffer
	wantHeader := "Example header"
	wantContent := "Example content"

	buffer.WriteString(wantHeader)
	buffer.WriteString("\n\n")
	buffer.WriteString(wantContent)
	buffer.WriteString("\n")

	gotHeader, gotContent, err := ReadSection(&buffer, "testfile")
	if err != nil {
		t.Error(err)
	}
	if gotHeader != wantHeader {
		t.Errorf("got header %q, want header %q", gotHeader, wantHeader)
	}

	if gotContent != wantContent {
		t.Errorf("got content %q, want content %q", gotContent, wantContent)
	}
}

func TestReadSectionTooLong(t *testing.T) {
	var buffer bytes.Buffer
	wantHeader := "Example header"
	wantContent := "Example content part 1\n\nExample content part 2"

	buffer.WriteString(wantHeader)
	buffer.WriteString("\n\n")
	buffer.WriteString("Example content part 1")
	buffer.WriteString("\n\n")
	buffer.WriteString("Example content part 2")
	buffer.WriteString("\n")

	gotHeader, gotContent, err := ReadSection(&buffer, "testfile")
	if err != nil {
		t.Error(err)
	}
	if gotHeader != wantHeader {
		t.Errorf("got header %q, want header %q", gotHeader, wantHeader)
	}

	if gotContent != wantContent {
		t.Errorf("got content %q, want content %q", gotContent, wantContent)
	}
}
