package models

import "testing"

func TestInferParentName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		ok       bool
	}{
		{
			name:     "species with aff",
			input:    "Proasellus aff. escolai",
			expected: "Proasellus",
			ok:       true,
		},
		{
			name:     "species with cf",
			input:    "Proasellus cf. escolai",
			expected: "Proasellus",
			ok:       true,
		},
		{
			name:     "genus with sp",
			input:    "Proasellus sp.",
			expected: "Proasellus",
			ok:       true,
		},
		{
			name:     "genus with sp and comment",
			input:    "Proasellus sp. (Styloss2)",
			expected: "Proasellus",
			ok:       true,
		},
		{
			name:     "subspecies",
			input:    "Proasellus anophtalmus resavicae",
			expected: "Proasellus anophtalmus",
			ok:       true,
		},
		{
			name:     "subgenus",
			input:    "Proasellus (coxalis)",
			expected: "Proasellus",
			ok:       true,
		},
		{
			name:     "subgenus ignored",
			input:    "Proasellus (coxalis) sp. (Styloss2)",
			expected: "Proasellus",
			ok:       true,
		},
		{
			name:     "species",
			input:    "Proasellus escolai",
			expected: "Proasellus",
			ok:       true,
		},
		{
			name:     "already genus",
			input:    "Proasellus",
			expected: "",
			ok:       false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
			ok:       false,
		},
		{
			name:     "only qualifier",
			input:    "sp.",
			expected: "",
			ok:       false,
		},
		{
			name:     "authority after species",
			input:    "Proasellus anophtalmus resavicae (Sket, 1959)",
			expected: "Proasellus anophtalmus",
			ok:       true,
		},
		{
			name:     "multiple spaces",
			input:    "Proasellus   aff.   escolai",
			expected: "Proasellus",
			ok:       true,
		},
		{
			name:     "species with spX",
			input:    "Proasellus sp1. (Maglia)",
			expected: "Proasellus",
			ok:       true,
		},
		{
			name:     "species with gen. n.",
			input:    "Candonidae gen. n. sp. I1, 2006",
			expected: "Candonidae gen. n.",
			ok:       true,
		},
		{
			name:     "species with gen. X",
			input:    "Candonidae gen. I1 sp. I1, 2006",
			expected: "Candonidae gen. I1",
			ok:       true,
		},
		{
			name:     "genus with gen. X",
			input:    "Candonidae gen. I1",
			expected: "Candonidae",
			ok:       true,
		},
		{
			name:     "species with n. sp.",
			input:    "Schellencandona n. sp. J4",
			expected: "Schellencandona",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := InferParentName(tt.input)

			if ok != tt.ok {
				t.Fatalf("InferParentName(%q) ok = %v, want %v", tt.input, ok, tt.ok)
			}

			if got != tt.expected {
				t.Errorf(
					"InferParentName(%q) = %q, want %q",
					tt.input,
					got,
					tt.expected,
				)
			}
		})
	}
}
