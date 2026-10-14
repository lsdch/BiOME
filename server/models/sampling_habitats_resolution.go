package models

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

// SamplingHabitatResolutionInput replaces the entire decision for one source value.
// Creations remain local to the batch until materialization.
type SamplingHabitatResolutionInput struct {
	InputText      string                `json:"input_text" minLength:"1"`
	Status         VocabResolutionStatus `json:"status"`
	HabitatIDs     []uuid.UUID           `json:"habitat_ids,omitempty"`
	CreateHabitats []HabitatInput        `json:"create_habitats,omitempty"`
}

type SamplingHabitatResolution struct {
	ImportID uuid.UUID `json:"import_id"`
	SamplingHabitatResolutionInput
}

func (i SamplingHabitatResolutionInput) Validate() error {
	if strings.TrimSpace(i.InputText) == "" {
		return WrapErrorPath(fmt.Errorf("source value is required"), "input_text")
	}
	switch i.Status {
	case biomedb.VocabResolutionStatusPending, biomedb.VocabResolutionStatusDiscard:
		if len(i.HabitatIDs) != 0 || len(i.CreateHabitats) != 0 {
			return WrapErrorPath(fmt.Errorf("pending and discarded values cannot have targets"), "status")
		}
	case biomedb.VocabResolutionStatusSelected:
		if len(i.HabitatIDs) == 0 || len(i.CreateHabitats) != 0 {
			return WrapErrorPath(fmt.Errorf("selected requires existing habitats and no creations"), "status")
		}
	case biomedb.VocabResolutionStatusRequestCreation:
		if len(i.CreateHabitats) == 0 {
			return WrapErrorPath(fmt.Errorf("request_creation requires at least one new habitat"), "create_habitats")
		}
	default:
		return WrapErrorPath(fmt.Errorf("unsupported manual resolution status %q", i.Status), "status")
	}
	ids := make(map[uuid.UUID]bool, len(i.HabitatIDs))
	for _, id := range i.HabitatIDs {
		if id == uuid.Nil || ids[id] {
			return WrapErrorPath(fmt.Errorf("habitat IDs must be nonzero and unique"), "habitat_ids")
		}
		ids[id] = true
	}
	for n, habitat := range i.CreateHabitats {
		if strings.TrimSpace(habitat.Name) == "" {
			return WrapErrorPath(fmt.Errorf("habitat name is required"), fmt.Sprintf("create_habitats.%d.name", n))
		}
		if habitat.Facet.IsSet && strings.TrimSpace(habitat.Facet.Value) == "" {
			return WrapErrorPath(fmt.Errorf("facet cannot be empty"), fmt.Sprintf("create_habitats.%d.facet", n))
		}
		if habitat.Description.IsSet && utf8.RuneCountInString(habitat.Description.Value) > 4000 {
			return WrapErrorPath(fmt.Errorf("description exceeds 4000 characters"), fmt.Sprintf("create_habitats.%d.description", n))
		}
		for _, previous := range i.CreateHabitats[:n] {
			if strings.EqualFold(previous.Name, habitat.Name) {
				return WrapErrorPath(fmt.Errorf("duplicate creation name %q", habitat.Name), "create_habitats")
			}
		}
	}
	return nil
}
