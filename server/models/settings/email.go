package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/edgedb/edgedb-go"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v3"
)

type EmailSettingsInput struct {
	FromName    string `edgedb:"from_name" json:"from_name" fake:"{firstname} {lastname}"`
	FromAddress string `edgedb:"from_address" json:"from_address" format:"email" fake:"{email}"`
	Host        string `edgedb:"host" json:"host" doc:"SMTP domain that handles email sending" format:"hostname" fake:"{domainname}"`
	Port        int32  `edgedb:"port" json:"port" doc:"SMTP port" minimum:"1" fake:"{number:10,99999}"`
	User        string `edgedb:"user" json:"user" doc:"SMTP login" format:"uri" fake:"{username}"`
	Password    string `edgedb:"password" json:"password" doc:"SMTP password" fake:"{password:true,true,true,true,true,20}"`
}

type EmailSettings struct {
	edgedb.Optional
	ID                 edgedb.UUID `edgedb:"id" json:"-"`
	EmailSettingsInput `edgedb:"$inline" json:",inline"`
}

func (e EmailSettingsInput) FromHeader() string {
	return fmt.Sprintf("%s <%s>", e.FromName, e.FromAddress)
}

func (e *EmailSettingsInput) Save(db edgedb.Executor) (*EmailSettings, error) {
	jsonData, _ := json.Marshal(e)
	var emailSettings EmailSettings
	if err := db.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			settings := (select admin::EmailSettings limit 1)
			select (if exists settings
				then (update admin::EmailSettings set {
					from_name := <str>data['from_name'],
					from_address := <str>data['from_address'],
					host := <str>data['host'],
					user := <str>data['user'],
					password := <str>data['password'],
					port := <int32>data['port'],
				})
				else (insert admin::EmailSettings {
					from_name := <str>data['from_name'],
					from_address := <str>data['from_address'],
					host := <str>data['host'],
					user := <str>data['user'],
					password := <str>data['password'],
					port := <int32>data['port'],
				})
			) { * } limit 1
		`,
		&emailSettings, jsonData,
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
