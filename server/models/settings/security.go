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
	MinPasswordStrength       int32  `edgedb:"min_password_strength" json:"min_password_strength" faker:"boundary_start=3, boundary_end=5"`
	AuthTokenLifetimeSeconds  int32  `edgedb:"auth_token_lifetime" json:"auth_token_lifetime" faker:"boundary_start=1, boundary_end=100000"`
	AccountTokenLifetimeHours int32  `edgedb:"account_token_lifetime" json:"account_token_lifetime" faker:"boundary_start=1, boundary_end=1000"`
	SecretKey                 string `edgedb:"jwt_secret_key" json:"jwt_secret_key" faker:"password,len=32"`
}

type SecuritySettings struct {
	ID                    edgedb.UUID `edgedb:"id" json:"id"`
	SecuritySettingsInput `edgedb:"$inline" json:",inline"`
}

func (s SecuritySettings) AuthTokenDuration() time.Duration {
	d, err := time.ParseDuration(
		fmt.Sprintf("%ds", s.AuthTokenLifetimeSeconds),
	)
	if err != nil {
		logrus.Fatalf("Failed to parse auth token duration: %v", err)
	}
	return d
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
	var settings SecuritySettings
	if err := db.QuerySingle(context.Background(),
		`with data := <json>$0
			select (update admin::SecuritySettings set {
				min_password_strength := <int16>data['min_password_strength'],
				auth_token_lifetime := <int16>data['auth_token_lifetime'],
				account_token_lifetime := <int16>data['account_token_lifetime'],
				jwt_secret_key := <str>data['jwt_secret_key']
			}) { * } limit 1`, &settings, jsonData,
	); err != nil {
		return nil, err
	}
	return &settings, nil
}
