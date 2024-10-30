package events

import "darco/proto/models/taxonomy"

type Spotting struct {
	TargetTaxa []taxonomy.TaxonInner `edgedb:"target_taxa" json:"target_taxa"`
}
