package settings

import (
	"context"
	"encoding/json"
	"os"

	"github.com/edgedb/edgedb-go"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v3"
)

type EmailSettingsInput struct {
	Host     string `edgedb:"host" json:"host" doc:"SMTP domain that handles email sending" format:"hostname" fake:"{domainname}"`
	Port     int32  `edgedb:"port" json:"port" doc:"SMTP port" minimum:"1" fake:"{number:10,99999}"`
	User     string `edgedb:"user" json:"user" doc:"SMTP login" format:"uri" fake:"{username}"`
	Password string `edgedb:"password" json:"password" doc:"SMTP password" fake:"{password:true,true,true,true,true,20}"`
}

type EmailSettings struct {
	edgedb.Optional
	ID                 edgedb.UUID `edgedb:"id" json:"-"`
	EmailSettingsInput `edgedb:"$inline" json:",inline"`
}

func (e *EmailSettingsInput) Save(db edgedb.Executor) (*EmailSettings, error) {
	jsonData, _ := json.Marshal(e)
	var emailSettings EmailSettings
	if err := db.QuerySingle(context.Background(),
		`with data := <json>$0,
			settings := (select admin::EmailSettings limit 1)
			select (if exists settings
				then (update admin::EmailSettings set {
					host := <str>data['host'],
					user := <str>data['user'],
					password := <str>data['password'],
					port := <int32>data['port'],
				})
				else (insert admin::EmailSettings {
					host := <str>data['host'],
					user := <str>data['user'],
					password := <str>data['password'],
					port := <int32>data['port'],
				})
			) { * } limit 1`,
		&emailSettings,
		jsonData,
	); err != nil {
		return nil, err
	}
	return &emailSettings, nil
}

func (e *EmailSettingsInput) Dialer() *gomail.Dialer {
	return gomail.NewDialer(e.Host, int(e.Port), e.User, e.Password)
}

func (e *EmailSettingsInput) TestConnection() error {
	_, err := e.Dialer().Dial()
	return err
}

func (e *EmailSettingsInput) WriteYAML(path string) error {
	yamlCfg, _ := yaml.Marshal(e)
	return os.WriteFile(path, yamlCfg, 0644)
}
