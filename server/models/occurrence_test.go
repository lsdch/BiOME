package models

import "testing"

func TestIdentification_String(t *testing.T) {
	tests := []struct {
		name           string
		identification Identification
		want           string
	}{
		{
			name: "taxon name only",
			identification: Identification{
				Taxon: Taxon{
					Name: "Asellus aquaticus",
				},
			},
			want: "Asellus aquaticus",
		},
		{
			name: "confer",
			identification: Identification{
				Taxon: Taxon{
					Name: "Asellus aquaticus",
				},
				Confer: true,
			},
			want: "Asellus cf. aquaticus",
		},
		{
			name: "addendum",
			identification: Identification{
				Taxon: Taxon{
					Name: "Asellus aquaticus",
				},
				Addendum: Optional[string]{
					Value: "juvenile",
					IsSet: true,
				},
			},
			want: "Asellus aquaticus juvenile",
		},
		{
			name: "confer and addendum",
			identification: Identification{
				Taxon: Taxon{
					Name: "Asellus aquaticus",
				},
				Confer: true,
				Addendum: Optional[string]{
					Value: "juvenile",
					IsSet: true,
				},
			},
			want: "Asellus cf. aquaticus juvenile",
		},
		{
			name: "multi-part taxon name",
			identification: Identification{
				Taxon: Taxon{
					Name: "Asellus aquaticus cavernicolus",
				},
				Confer: true,
			},
			want: "Asellus aquaticus cf. cavernicolus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.identification.String()

			if got != tt.want {
				t.Errorf("Identification.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
