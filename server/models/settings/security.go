package settings

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type SecuritySettingsInput struct {
	MinPasswordStrength         int32 `edgedb:"min_password_strength" json:"min_password_strength" doc:"The level of complexity required for account passwords." minimum:"3" maximum:"5" fake:"{number:3,5}"`
	RefreshTokenLifetimeHours   int32 `edgedb:"refresh_token_lifetime" json:"refresh_token_lifetime" doc:"User session lifetime in hours" minimum:"1" fake:"{number:1,2}"`
	InvitationTokenLifetimeDays int32 `edgedb:"invitation_token_lifetime" json:"invitation_token_lifetime" doc:"Invitation token lifetime in days" minimum:"1" fake:"{number:1,7}"`
}

type SecuritySettings struct {
	ID                    edgedb.UUID `edgedb:"id" json:"-"`
	SecuritySettingsInput `edgedb:"$inline" json:",inline"`
	SecretKey             string `edgedb:"jwt_secret_key" json:"-" doc:"Used to verify session tokens. Changing it will revoke all currently active user sessions." minLength:"32" fake:"{password:true,true,true,true,true,32}"`
}

func (s SecuritySettings) RefreshTokenDuration() time.Duration {
	d, err := time.ParseDuration(
		fmt.Sprintf("%dh", s.RefreshTokenLifetimeHours),
	)
	if err != nil {
		logrus.Fatalf("Failed to parse token duration: %v", err)
	}
	return d
}

func (s SecuritySettings) InvitationTokenDuration() time.Duration {
	d, err := time.ParseDuration(
		fmt.Sprintf("%dd", s.InvitationTokenLifetimeDays),
	)
	if err != nil {
		logrus.Fatalf("Failed to parse token duration: %v", err)
	}
	return d
}

func generateSecretKeyJWT() string {
	// Determine the number of bytes needed for the given key length
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		logrus.Fatalf("Failed to generate a 32 bytes secret key for JWT encryption.")
	}
	// Encode the random bytes to base64 to create a printable string
	encodedKey := base64.URLEncoding.EncodeToString(key)
	return encodedKey
}

func (s *SecuritySettings) RefreshSecretKey(db edgedb.Executor) error {
	secretKey := generateSecretKeyJWT()
	if err := db.Execute(context.Background(),
		`update admin::SecuritySettings set { jwt_secret_key := <str>$0 }`,
		secretKey,
	); err != nil {
		return err
	}

	s.SecretKey = secretKey
	return nil
}

func (input *SecuritySettingsInput) Save(db edgedb.Executor) (*SecuritySettings, error) {
	jsonData, _ := json.Marshal(input)
	var s SecuritySettings
	if err := db.QuerySingle(context.Background(),
		`with data := <json>$0
			select (update admin::SecuritySettings set {
				min_password_strength := <int32>data['min_password_strength'],
				refresh_token_lifetime := <int32>data['refresh_token_lifetime'],
				invitation_token_lifetime := <int32>data['invitation_token_lifetime'],
			}) { * } limit 1`, &s, jsonData,
	); err != nil {
		return nil, err
	}
	return &s, nil
}
