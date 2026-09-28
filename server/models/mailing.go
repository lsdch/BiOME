package models

import (
	"time"

	"github.com/lsdch/biome/db/biomedb"
)

type Mailing struct {
	Enabled         bool      `json:"enabled"`
	MailFromAddress string    `json:"mail_from_address"`
	MailFromName    string    `json:"mail_from_name"`
	SmtpHost        *string   `json:"smtp_host"`
	SmtpPort        *int32    `json:"smtp_port"`
	SmtpUser        *string   `json:"smtp_user"`
	SmtpPassword    *string   `json:"smtp_password"`
	LastUpdated     time.Time `json:"last_updated"`
}

func (m Mailing) IsConfigured() bool {
	return m.SmtpHost != nil &&
		*m.SmtpHost != "" &&
		m.SmtpPort != nil &&
		*m.SmtpPort > 0 &&
		*m.SmtpPort <= 65535 &&
		m.MailFromAddress != "" &&
		m.MailFromName != "" &&
		m.LastUpdated.After(time.Time{})
}

func MailingFromDB(m biomedb.Mailing) Mailing {
	return Mailing{
		Enabled:         m.Enabled,
		MailFromAddress: m.MailFromAddress,
		MailFromName:    m.MailFromName,
		SmtpHost:        m.SmtpHost,
		SmtpPort:        m.SmtpPort,
		SmtpUser:        m.SmtpUser,
		SmtpPassword:    m.SmtpPassword,
		LastUpdated:     m.LastUpdated,
	}
}

type UpsertMailingParams struct {
	MailFromAddress string  `json:"mail_from_address"`
	MailFromName    string  `json:"mail_from_name"`
	SmtpHost        *string `json:"smtp_host"`
	SmtpPort        *int32  `json:"smtp_port"`
	SmtpUser        *string `json:"smtp_user"`
	SmtpPassword    *string `json:"smtp_password"`
}

func (p UpsertMailingParams) ToDBParams() biomedb.UpsertMailingParams {
	return biomedb.UpsertMailingParams{
		MailFromAddress: p.MailFromAddress,
		MailFromName:    p.MailFromName,
		SmtpHost:        p.SmtpHost,
		SmtpPort:        p.SmtpPort,
		SmtpUser:        p.SmtpUser,
		SmtpPassword:    p.SmtpPassword,
	}
}
