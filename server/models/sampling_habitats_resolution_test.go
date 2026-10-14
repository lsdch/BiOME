package models

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/stretchr/testify/require"
)

func TestSamplingHabitatResolutionValidate(t *testing.T) {
	id := uuid.New()
	creation := HabitatInput{Name: "Cave", Facet: NewOptional("New facet")}
	for _, tc := range []struct {
		name      string
		status    VocabResolutionStatus
		ids       []uuid.UUID
		creations []HabitatInput
		valid     bool
	}{
		{name: "pending", status: biomedb.VocabResolutionStatusPending, valid: true},
		{name: "discard", status: biomedb.VocabResolutionStatusDiscard, valid: true},
		{name: "multiple existing", status: biomedb.VocabResolutionStatusSelected, ids: []uuid.UUID{id, uuid.New()}, valid: true},
		{name: "creation with new facet", status: biomedb.VocabResolutionStatusRequestCreation, creations: []HabitatInput{creation}, valid: true},
		{name: "mixed targets", status: biomedb.VocabResolutionStatusRequestCreation, ids: []uuid.UUID{id}, creations: []HabitatInput{creation}, valid: true},
		{name: "empty selected", status: biomedb.VocabResolutionStatusSelected},
		{name: "empty creation", status: biomedb.VocabResolutionStatusRequestCreation},
		{name: "discard with targets", status: biomedb.VocabResolutionStatusDiscard, ids: []uuid.UUID{id}},
		{name: "selected with creation", status: biomedb.VocabResolutionStatusSelected, ids: []uuid.UUID{id}, creations: []HabitatInput{creation}},
		{name: "automatic status cannot be submitted", status: biomedb.VocabResolutionStatusAutoResolved, ids: []uuid.UUID{id}},
		{name: "invalid status", status: "invalid"},
		{name: "nil UUID", status: biomedb.VocabResolutionStatusSelected, ids: []uuid.UUID{uuid.Nil}},
		{name: "duplicate ID", status: biomedb.VocabResolutionStatusSelected, ids: []uuid.UUID{id, id}},
		{name: "duplicate name", status: biomedb.VocabResolutionStatusRequestCreation, creations: []HabitatInput{creation, {Name: "cAVE"}}},
		{name: "blank name", status: biomedb.VocabResolutionStatusRequestCreation, creations: []HabitatInput{{Name: " "}}},
		{name: "blank facet", status: biomedb.VocabResolutionStatusRequestCreation, creations: []HabitatInput{{Name: "Cave", Facet: NewOptional(" ")}}},
		{name: "long description", status: biomedb.VocabResolutionStatusRequestCreation, creations: []HabitatInput{{Name: "Cave", Description: NewOptional(strings.Repeat("é", 4001))}}},
		{name: "unicode description at limit", status: biomedb.VocabResolutionStatusRequestCreation, creations: []HabitatInput{{Name: "Cave", Description: NewOptional(strings.Repeat("é", 4000))}}, valid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input := SamplingHabitatResolutionInput{InputText: "source", Status: tc.status, HabitatIDs: tc.ids, CreateHabitats: tc.creations}
			if tc.valid {
				require.NoError(t, input.Validate())
			} else {
				require.Error(t, input.Validate())
			}
		})
	}
	require.Error(t, (SamplingHabitatResolutionInput{InputText: " ", Status: biomedb.VocabResolutionStatusDiscard}).Validate())
}
