package readers

import (
	"bytes"
	"testing"
)

func TestReadSection(t *testing.T) {
	testCases := []struct {
		input       *bytes.Buffer
		wantHeader  string
		wantContent string
	}{
		{
			input:       bytes.NewBufferString("Example header\n\nExample content\n"),
			wantHeader:  "Example header",
			wantContent: "Example content",
		},
		{
			input:       bytes.NewBufferString("Example header\n\nExample content\npart 1\n\nExample content part 2\n\nExample content part 3\n"),
			wantHeader:  "Example header",
			wantContent: "Example content part 1\n\nExample content part 2\n\nExample content part 3",
		},
	}

	for _, tc := range testCases {
		gotHeader, gotContent, err := ReadSection(tc.input)
		if err != nil {
			t.Error(err)
		}
		if gotHeader != tc.wantHeader {
			t.Errorf("got header %q, want header %q", gotHeader, tc.wantHeader)
		}

		if gotContent != tc.wantContent {
			t.Errorf("got content %q, want content %q", gotContent, tc.wantContent)
		}
	}
}
