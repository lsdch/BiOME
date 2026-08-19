package models

import (
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type Habitat struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description Optional[string] `json:"description,omitzero"`
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
	Description Optional[string]    `json:"description,omitzero"`
	ParentID    Optional[uuid.UUID] `json:"parent_id,omitzero"`
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
	Description Optional[string] `json:"description,omitzero" doc:"Optional habitat description"`
}

type HabitatGroupInput struct {
	Label     string           `json:"label" doc:"Name for the group of habitat tags" example:"Water flow" minLength:"3" maxLength:"32"`
	Depends   Optional[string] `json:"depends,omitzero" doc:"Habitat tag that this group is a refinement of" example:"Aquatic, Surface"`
	Exclusive Optional[bool]   `json:"exclusive_elements,omitzero"`
	Elements  []HabitatInput   `json:"elements" minItems:"1"`
}

func (i HabitatGroupInput) ToDBParams() biomedb.InsertHabitatGroupParams {
	return biomedb.InsertHabitatGroupParams{
		Label:             i.Label,
		Depends:           i.Depends.ToPtr(),
		ExclusiveElements: i.Exclusive.Value,
	}
}

type HabitatGroupWithElements struct {
	HabitatGroup
	Elements []Habitat `json:"elements"`
}

type HabitatUpdate struct {
	Label       Optional[string]     `gel:"label" json:"label,omitempty"`
	Description OptionalNull[string] `gel:"description" json:"description,omitempty"`
}

func (u HabitatUpdate) ToDBParams(id uuid.UUID) biomedb.UpdateHabitatParams {
	return biomedb.UpdateHabitatParams{
		HabitatID:      id,
		Label:          u.Label.ToPtr(),
		SetDescription: u.Description.IsSet,
		Description:    u.Description.ToPtr(),
	}
}

type HabitatGroupUpdate struct {
	Label          Optional[string]            `gel:"label" json:"label,omitempty"`
	Depends        OptionalNull[uuid.UUID]     `gel:"depends" json:"depends,omitempty"`
	Exclusive      Optional[bool]              `gel:"exclusive_elements" json:"exclusive_elements,omitempty"`
	CreateElements []HabitatInput              `json:"create_elements,omitempty"`
	UpdateElements map[uuid.UUID]HabitatUpdate `json:"update_elements,omitempty"`
	DeleteElements []uuid.UUID                 `json:"delete_elements,omitempty"`
}

func (u HabitatGroupUpdate) HasUpdateInfos() bool {
	return u.Label.IsSet || u.Depends.IsSet || u.Exclusive.IsSet
}

func (u HabitatGroupUpdate) ToDBParams(id uuid.UUID) biomedb.UpdateHabitatGroupInfoParams {
	return biomedb.UpdateHabitatGroupInfoParams{
		GroupID:            id,
		Label:              u.Label.ToPtr(),
		ParentHabitatID:    UUIDOpt(u.Depends),
		SetParentHabitatID: u.Depends.IsSet,
		ExclusiveElements:  u.Exclusive.ToPtr(),
	}
}
