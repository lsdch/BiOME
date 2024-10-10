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
	MinPasswordStrength       int32  `edgedb:"min_password_strength" json:"min_password_strength" doc:"The level of complexity required for account passwords." minimum:"3" maximum:"5" fake:"{number:3,5}"`
	AuthTokenLifetimeMinutes  int32  `edgedb:"auth_token_lifetime" json:"auth_token_lifetime" doc:"JWT session lifetime in minutes" minimum:"15" fake:"{number:15,1000}"`
	RefreshTokenLifetimeDays  int32  `edgedb:"refresh_token_lifetime" json:"refresh_token_lifetime" doc:"User session lifetime in days" minimum:"1" fake:"{number:1,2}"`
	AccountTokenLifetimeHours int32  `edgedb:"account_token_lifetime" json:"account_token_lifetime" minimum:"2" doc:"Account manipulation token lifetime in hours" fake:"{number:2,48}"`
	SecretKey                 string `edgedb:"jwt_secret_key" json:"jwt_secret_key" doc:"Used to verify session tokens. Changing it will revoke all currently active user sessions." minLength:"32" fake:"{password:true,true,true,true,true,32}"`
}

type SecuritySettings struct {
	ID                    edgedb.UUID `edgedb:"id" json:"-"`
	SecuritySettingsInput `edgedb:"$inline" json:",inline"`
}

func (s SecuritySettings) AuthTokenDuration() time.Duration {
	d, err := time.ParseDuration(
		fmt.Sprintf("%dm", s.AuthTokenLifetimeMinutes),
	)
	if err != nil {
		logrus.Fatalf("Failed to parse auth token duration: %v", err)
	}
	return d
}

func (s SecuritySettings) CookieMaxAge() int {
	return int(s.AuthTokenDuration()) * 60
}

func (s SecuritySettings) AccountTokenDuration() time.Duration {
	d, err := time.ParseDuration(
		fmt.Sprintf("%dh", s.AccountTokenLifetimeHours),
	)
	if err != nil {
		logrus.Fatalf("Failed to parse account token duration: %v", err)
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

func (input *SecuritySettingsInput) Save(db edgedb.Executor) (*SecuritySettings, error) {
	jsonData, _ := json.Marshal(input)
	var s SecuritySettings
	if err := db.QuerySingle(context.Background(),
		`with data := <json>$0
			select (update admin::SecuritySettings set {
				min_password_strength := <int32>data['min_password_strength'],
				auth_token_lifetime := <int32>data['auth_token_lifetime'],
				refresh_token_lifetime := <int32>data['refresh_token_lifetime'],
				account_token_lifetime := <int32>data['account_token_lifetime'],
				jwt_secret_key := <str>data['jwt_secret_key']
			}) { * } limit 1`, &s, jsonData,
	); err != nil {
		return nil, err
	}
	return &s, nil
}
