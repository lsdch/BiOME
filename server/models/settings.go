package models

import "github.com/lsdch/biome/db/biomedb"

type InstanceSettings struct {
	Title                  string           `json:"title"`
	Subtitle               Optional[string] `json:"subtitle,omitempty"`
	Description            Optional[string] `json:"description,omitempty"`
	AdminEmail             string           `json:"admin_email"`
	IsPublic               bool             `json:"is_public"`
	AccountRequestsEnabled bool             `json:"account_requests_enabled"`
	MailFromAddress        string           `json:"mail_from_address"`
	MailFromName           string           `json:"mail_from_name"`
	MolecularDataEnabled   bool             `json:"molecular_data_enabled"`
}

func SettingsFromDB(s biomedb.Setting) InstanceSettings {
	return InstanceSettings{
		Title:                  s.AppName,
		Subtitle:               NewOptionalFromPtr(s.AppSubtitle),
		Description:            NewOptionalFromPtr(s.AppDescription),
		AdminEmail:             s.AdminEmail,
		IsPublic:               s.IsPublic,
		AccountRequestsEnabled: s.AccountRequestsEnabled,
		MailFromAddress:        s.MailFromAddress,
		MailFromName:           s.MailFromName,
		MolecularDataEnabled:   s.MolecularDataEnabled,
	}
}

type InstanceSettingsUpdate struct {
	Title           Optional[string]     `json:"title,omitempty"`
	Subtitle        OptionalNull[string] `json:"subtitle,omitempty"`
	Description     OptionalNull[string] `json:"description,omitempty"`
	AdminEmail      Optional[string]     `json:"admin_email,omitempty"`
	MailFromAddress Optional[string]     `json:"mail_from_address,omitempty"`
	MailFromName    Optional[string]     `json:"mail_from_name,omitempty"`
}

func (s *InstanceSettingsUpdate) ToParams() biomedb.UpdateInstanceSettingsParams {
	return biomedb.UpdateInstanceSettingsParams{
		AppName:           s.Title.ToPtr(),
		SetAppSubtitle:    s.Subtitle.IsSet,
		AppSubtitle:       s.Subtitle.ToPtr(),
		SetAppDescription: s.Description.IsSet,
		AppDescription:    s.Description.ToPtr(),
		AdminEmail:        s.AdminEmail.ToPtr(),
		MailFromAddress:   s.MailFromAddress.ToPtr(),
		MailFromName:      s.MailFromName.ToPtr(),
	}
}
