package formatters_test

import (
	"testing"

	"github.com/bartekpacia/database-tools/formatters"
)

func TestToContent(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{
			input: "Example content\n",
			want:  "Example content",
		},
		{
			input: "Example\ncontent part 1\n\nExample content\npart 2\n",
			want:  "Example content part 1\n\nExample content part 2",
		},
	}

	for _, tc := range testCases {
		got := formatters.ToContent(tc.input)
		if got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}

		if got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}

func TestToSection(t *testing.T) {
	testCases := []struct {
		input       string
		wantHeader  string
		wantContent string
	}{
		{
			input:       "Example header\n\nExample content\n",
			wantHeader:  "Example header",
			wantContent: "Example content",
		},
		{
			input:       "Example header\n\nExample content\npart 1\n\nExample content part 2\n\nExample content part 3\n",
			wantHeader:  "Example header",
			wantContent: "Example content part 1\n\nExample content part 2\n\nExample content part 3",
		},
	}

	for _, tc := range testCases {
		gotHeader, gotContent := formatters.ToSection(tc.input)
		if gotHeader != tc.wantHeader {
			t.Errorf("got header %q, want header %q", gotHeader, tc.wantHeader)
		}

		if gotContent != tc.wantContent {
			t.Errorf("got content %q, want content %q", gotContent, tc.wantContent)
		}
	}
}
