package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"testing"
	"time"
)

func TestEmailConfirmationToken(t *testing.T) {
	client := db.Client()
	p, err := FakePendingUserInput(t).Register(client)
	if err != nil {
		t.Fatalf("%v", err)
	}
	token, err := p.User.CreateAccountToken(client, people.EmailConfirmationToken)
	if err != nil {
		t.Fatalf("Token generation failed: %v", err)
	}
	if len(token) == 0 {
		t.Fatalf("Empty token generated")
	}

	time.Sleep(time.Second)
	user, ok := people.ValidateAccountToken(client, token, people.EmailConfirmationToken)
	if !ok {
		t.Fatalf("Failed to validate account token")
	}
	if !user.EmailConfirmed {
		t.Fatalf("User email was not confirmed")
	}
}
