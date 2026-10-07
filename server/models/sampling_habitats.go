package models

import (
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"time"
)

type Habitat struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Facet       Optional[string] `json:"facet,omitzero"`
	Description Optional[string] `json:"description,omitzero"`
	CreatedAt   time.Time        `json:"created_at"`
}

func HabitatFromDB(h biomedb.Habitat) Habitat {
	return Habitat{ID: h.ID, Name: h.Label, Facet: NewOptionalFromPtr(h.Facet), Description: NewOptionalFromPtr(h.Description), CreatedAt: h.CreatedAt}
}

type HabitatInput struct {
	Name        string           `json:"name" yaml:"name" doc:"Habitat name" minLength:"1"`
	Facet       Optional[string] `json:"facet,omitzero" yaml:"facet" doc:"Optional facet name" minLength:"1"`
	Description Optional[string] `json:"description,omitzero" yaml:"description" maxLength:"4000"`
}

func (i HabitatInput) ToDBParams() biomedb.InsertHabitatParams {
	return biomedb.InsertHabitatParams{Label: i.Name, Facet: i.Facet.ToPtr(), Description: i.Description.ToPtr()}
}

type HabitatUpdate struct {
	Name        Optional[string]     `json:"name,omitempty" minLength:"1"`
	Facet       OptionalNull[string] `json:"facet,omitempty" minLength:"1"`
	Description OptionalNull[string] `json:"description,omitempty" maxLength:"4000"`
}

func (u HabitatUpdate) ToDBParams(id uuid.UUID) biomedb.UpdateHabitatParams {
	return biomedb.UpdateHabitatParams{HabitatID: id, Label: u.Name.ToPtr(), SetFacet: u.Facet.IsSet, Facet: u.Facet.ToPtr(), SetDescription: u.Description.IsSet, Description: u.Description.ToPtr()}
}
