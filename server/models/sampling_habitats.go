package models

import (
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type Habitat struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description Optional[string] `json:"description,omitempty"`
}

func HabitatFromDB(h biomedb.Habitat) Habitat {
	return Habitat{
		ID:          h.ID,
		Name:        h.Label,
		Description: NewOptionalFromPtr(h.Description),
	}
}

type HabitatWithGroupName struct {
	Habitat
	GroupName string `json:"group_name"`
}

func HabitatWithGroupNameFromDB(h biomedb.Habitat, groupName string) HabitatWithGroupName {
	return HabitatWithGroupName{
		Habitat:   HabitatFromDB(h),
		GroupName: groupName,
	}
}

type HabitatGroup struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name"`
	Description Optional[string]    `json:"description,omitempty"`
	ParentID    Optional[uuid.UUID] `json:"parent_id,omitempty"`
}

func HabitatGroupFromDB(h biomedb.HabitatGroup) HabitatGroup {
	return HabitatGroup{
		ID:          h.ID,
		Name:        h.Label,
		Description: NewOptionalFromPtr(h.Description),
		ParentID:    NewOptionalFromUUID(h.ParentHabitatID),
	}
}

type HabitatInput struct {
	Label       string           `json:"label" doc:"A short label for the habitat." minLength:"3" maxLength:"32" example:"Lotic"`
	Description Optional[string] `json:"description,omitempty" doc:"Optional habitat description"`
}

type HabitatGroupInput struct {
	Label       string              `json:"label" doc:"Name for the group of habitat tags" example:"Water flow" minLength:"3" maxLength:"32"`
	Description Optional[string]    `json:"description,omitempty" doc:"Optional description for the habitat group"`
	Depends     Optional[uuid.UUID] `json:"depends,omitempty" doc:"Habitat tag that this group is a refinement of" example:"Aquatic, Surface"`
	Exclusive   Optional[bool]      `json:"exclusive_elements,omitempty"`
	Elements    []HabitatInput      `json:"elements" minItems:"1"`
}

type HabitatGroupWithElements struct {
	HabitatGroup
	Elements []Habitat `json:"elements"`
}
