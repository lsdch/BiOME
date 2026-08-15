package bibliography

import (
	"reflect"
	"testing"

	"github.com/lsdch/biome/types"
)

func TestExtractAuthorString(t *testing.T) {
	tests := []struct {
		name     string
		verbatim string
		expected string
	}{
		{
			name:     "simple authors",
			verbatim: "Smith J., Doe A. 2020. Great paper.",
			expected: "Smith J., Doe A",
		},
		{
			name:     "year in parentheses",
			verbatim: "Smith J., Doe A. (2020). Great paper.",
			expected: "Smith J., Doe A",
		},
		{
			name:     "multiple spaces",
			verbatim: "Smith J.,   Doe A.    2020. Great paper.",
			expected: "Smith J.,   Doe A",
		},
		{
			name:     "no year",
			verbatim: "Smith J., Doe A. Great paper.",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractAuthorString(tt.verbatim)

			if got != tt.expected {
				t.Errorf("ExtractAuthorString() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestExtractYear(t *testing.T) {
	tests := []struct {
		name     string
		verbatim string
		expected *int32
	}{
		{
			name:     "simple year",
			verbatim: "Smith J. 2020. Great paper.",
			expected: ptrInt32(2020),
		},
		{
			name:     "year in parentheses",
			verbatim: "Smith J. (2020). Great paper.",
			expected: ptrInt32(2020),
		},
		{
			name:     "no year",
			verbatim: "Smith J. Great paper.",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractYear(tt.verbatim)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ExtractYear() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExtractDOI(t *testing.T) {
	tests := []struct {
		name     string
		verbatim string
		expected *types.DOI
	}{
		{
			name:     "doi prefix",
			verbatim: "Smith J. 2020. doi:10.1234/example.paper",
			expected: ptrDOI("10.1234/example.paper"),
		},
		{
			name:     "doi url",
			verbatim: "https://doi.org/10.1000/xyz123",
			expected: ptrDOI("10.1000/xyz123"),
		},
		{
			name:     "doi with punctuation",
			verbatim: "doi:10.5555/test.",
			expected: ptrDOI("10.5555/test"),
		},
		{
			name:     "no doi",
			verbatim: "Smith J. 2020. Great paper.",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractDOI(tt.verbatim)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ExtractDOI() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseAuthors(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected []string
	}{
		{
			name: "comma separated classic",
			raw:  "Smith J., Doe A.",
			expected: []string{
				"Smith J.",
				"Doe A.",
			},
		},
		{
			name: "semicolon separated",
			raw:  "Smith J.; Doe A.",
			expected: []string{
				"Smith J.",
				"Doe A.",
			},
		},
		{
			name: "and separator",
			raw:  "Smith J. and Doe A.",
			expected: []string{
				"Smith J.",
				"Doe A.",
			},
		},
		{
			name: "ampersand separator",
			raw:  "Smith J. & Doe A.",
			expected: []string{
				"Smith J.",
				"Doe A.",
			},
		},
		{
			name: "lastname comma initials",
			raw:  "Smith, J., Doe, A.",
			expected: []string{
				"Smith, J.",
				"Doe, A.",
			},
		},
		{
			name: "extra spaces",
			raw:  "  Smith J.   ,   Doe A. ",
			expected: []string{
				"Smith J.",
				"Doe A.",
			},
		},
		{
			name:     "empty",
			raw:      "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseAuthors(tt.raw)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseAuthors() = %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestNormalizeSpaces(t *testing.T) {
	got := normalizeSpaces("  Smith   J.   Doe ")
	expected := "Smith J. Doe"

	if got != expected {
		t.Errorf("normalizeSpaces() = %q, want %q", got, expected)
	}
}

func TestExtractAndParseAuthors(t *testing.T) {
	verbatim := "Wächtler W (1937) II. Ordnung: Isopoda, Asseln. In: Brohmer P, Ehrmann P, Ulmer G (Eds), Die Tierwelt Mitteleuropas. Quelle & Meyer, Leipzig, 225–317."

	raw := ExtractAuthorString(verbatim)
	got := ParseAuthors(raw)

	expected := []string{"Wächtler W"}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %#v, want %#v", got, expected)
	}
}

func ptrInt32(v int32) *int32 {
	return &v
}

func ptrDOI(v string) *types.DOI {
	doi := types.ParseDOI(v)
	return &doi
}
