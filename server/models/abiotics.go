package models

import (
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type AbioticParam struct {
	ID          uuid.UUID        `json:"id" format:"uuid"`
	Name        string           `json:"name"`
	Code        string           `json:"code"`
	Description Optional[string] `json:"description,omitempty"`
	Unit        string           `json:"unit"`
}

func AbioticParamFromDB(p biomedb.AbioticParam) AbioticParam {
	return AbioticParam{
		ID:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Description: NewOptionalFromPtr(p.Description),
		Unit:        p.Unit,
	}
}

type AbioticParamInput struct {
	Name        string           `json:"name" doc:"Name of the abiotic parameter." minLength:"3" maxLength:"64" example:"Temperature"`
	Code        string           `json:"code" doc:"Short code for the abiotic parameter." minLength:"2" maxLength:"16" example:"temp"`
	Description Optional[string] `json:"description,omitempty" doc:"Optional description of the abiotic parameter."`
	Unit        string           `json:"unit" doc:"Unit of measurement for the abiotic parameter." minLength:"1" maxLength:"16" example:"°C"`
}

func (i AbioticParamInput) ToDB() biomedb.CreateAbioticParamParams {
	return biomedb.CreateAbioticParamParams{
		Name:        i.Name,
		Code:        i.Code,
		Description: i.Description.ToPtr(),
		Unit:        i.Unit,
	}
}
