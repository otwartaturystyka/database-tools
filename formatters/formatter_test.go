package formatters_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opentouristics/database-tools/formatters"
)

func TestToContent(t *testing.T) {
	testCases := []struct {
		input map[string]string
		want  map[string]string
	}{
		{
			input: map[string]string{"en": "Example content\n"},
			want:  map[string]string{"en": "Example content"},
		},
		{
			input: map[string]string{"en": "Example\ncontent part 1\n\nExample content\npart 2\n"},
			want:  map[string]string{"en": "Example content part 1\n\nExample content part 2"},
		},
	}

	for _, tc := range testCases {
		got := formatters.ToContent(tc.input)

		if !cmp.Equal(got, tc.want) {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}

func TestToSection(t *testing.T) {
	testCases := []struct {
		input       map[string]string
		wantHeader  map[string]string
		wantContent map[string]string
	}{
		{
			input:       map[string]string{"en": "Example header\n\nExample content\n"},
			wantHeader:  map[string]string{"en": "Example header"},
			wantContent: map[string]string{"en": "Example content"},
		},
		{
			input:       map[string]string{"en": "Example header\n\nExample content\npart 1\n\nExample content part 2\n\nExample content part 3\n"},
			wantHeader:  map[string]string{"en": "Example header"},
			wantContent: map[string]string{"en": "Example content part 1\n\nExample content part 2\n\nExample content part 3"},
		},
	}

	for _, tc := range testCases {
		gotHeader, gotContent := formatters.ToSection(tc.input)
		if !cmp.Equal(gotHeader, tc.wantHeader) {
			t.Errorf("got header %q, want header %q", gotHeader, tc.wantHeader)
		}

		if !cmp.Equal(gotContent, tc.wantContent) {
			t.Errorf("got content %q, want content %q", gotContent, tc.wantContent)
		}
	}
}
