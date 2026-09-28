package models

import "github.com/lsdch/biome/db/biomedb"

type InstanceSettings struct {
	Title                  string           `json:"title"`
	Subtitle               Optional[string] `json:"subtitle,omitzero"`
	Description            Optional[string] `json:"description,omitzero"`
	AdminEmail             string           `json:"admin_email"`
	IsPublic               bool             `json:"is_public"`
	AccountRequestsEnabled bool             `json:"account_requests_enabled"`
	MolecularDataEnabled   bool             `json:"molecular_data_enabled"`
}

func SettingsFromDB(s biomedb.Setting) InstanceSettings {
	return InstanceSettings{
		Title:                  s.AppName,
		AdminEmail:             s.AdminEmail,
		Subtitle:               NewOptionalFromPtr(s.AppSubtitle),
		Description:            NewOptionalFromPtr(s.AppDescription),
		IsPublic:               s.IsPublic,
		AccountRequestsEnabled: s.AccountRequestsEnabled,
		MolecularDataEnabled:   s.MolecularDataEnabled,
	}
}

type InstanceSettingsUpdate struct {
	Title       Optional[string]     `json:"title,omitzero"`
	Subtitle    OptionalNull[string] `json:"subtitle,omitempty"`
	Description OptionalNull[string] `json:"description,omitempty"`
	AdminEmail  Optional[string]     `json:"admin_email,omitzero"`
}

func (s *InstanceSettingsUpdate) ToParams() biomedb.UpdateInstanceSettingsParams {
	return biomedb.UpdateInstanceSettingsParams{
		AppName:           s.Title.ToPtr(),
		AdminEmail:        s.AdminEmail.ToPtr(),
		SetAppSubtitle:    s.Subtitle.IsSet,
		AppSubtitle:       s.Subtitle.ToPtr(),
		SetAppDescription: s.Description.IsSet,
		AppDescription:    s.Description.ToPtr(),
	}
}
