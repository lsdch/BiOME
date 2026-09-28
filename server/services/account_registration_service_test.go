package services

import (
	"context"
	"encoding/base64"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/models"
	email_templates "github.com/lsdch/biome/templates"
	"github.com/stretchr/testify/require"
)

func TestInvitationTokenAndURL(t *testing.T) {
	token, err := GenerateSecureToken()
	require.NoError(t, err)
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	require.NoError(t, err)
	require.Len(t, decoded, 32)
	other, err := GenerateSecureToken()
	require.NoError(t, err)
	require.NotEqual(t, token, other)
	require.Equal(t, "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", HashToken("abc"))
	base, err := url.Parse("https://biome.example/app?lang=fr")
	require.NoError(t, err)
	service := NewRegistrationService(config.Config{AppPublicBaseURL: *base}, nil)
	link := service.InvitationURL("register", token)
	require.Equal(t, "/app/register", link.Path)
	require.Equal(t, token, link.Query().Get("token"))
	require.Equal(t, "fr", link.Query().Get("lang"))
	require.Empty(t, base.Query().Get("token"))
}

func TestCreateInvitationRejectsInvalidInputBeforeDatabase(t *testing.T) {
	service := NewRegistrationService(config.Config{}, nil)
	valid := models.CreateInvitationParams{Email: "alice@example.org", InviteeName: "Alice", Role: "contributor", ExpiresAt: time.Now().Add(time.Hour)}
	for _, mutate := range []func(*models.CreateInvitationParams){
		func(p *models.CreateInvitationParams) { p.Email = "invalid" },
		func(p *models.CreateInvitationParams) { p.Email = "Alice <alice@example.org>" },
		func(p *models.CreateInvitationParams) { p.InviteeName = " " },
		func(p *models.CreateInvitationParams) { p.Role = "invalid" },
		func(p *models.CreateInvitationParams) { p.ExpiresAt = time.Now().Add(-time.Hour) },
	} {
		input := valid
		mutate(&input)
		_, err := service.CreateInvitation(context.Background(), nil, input)
		require.ErrorIs(t, err, ErrInvalidRegistration)
	}
}

func TestInvitationEarlyFailures(t *testing.T) {
	service := NewRegistrationService(config.Config{}, nil)
	_, err := service.GetInvitationByToken(nil, context.Background(), "")
	require.ErrorIs(t, err, ErrInvalidInvitation)
	_, err = service.RegisterFromInvitation(nil, context.Background(), "", models.RegisterFromInvitationParams{})
	require.ErrorIs(t, err, ErrInvalidInvitation)
	require.ErrorIs(t, service.SendInvitationEmail(context.Background(), nil, models.InvitationResult{}, "register"), ErrMailerUnavailable)
	require.ErrorContains(t, service.SendInvitationEmail(context.Background(), nil, models.InvitationResult{}, "register"), "URL configuration")
}

func TestRegistrationHashRejectsLongPassword(t *testing.T) {
	for _, password := range []string{strings.Repeat("a", 73), strings.Repeat("é", 37)} {
		input := models.RegisterFromInvitationParams{Password: password}
		params, err := input.ToParams("hash")
		require.Error(t, err)
		require.Empty(t, params.PasswordHash)
	}
}

func TestInvitationEmailTemplate(t *testing.T) {
	message := "Welcome <script>alert(1)</script>"
	data := email_templates.InvitationData{
		Name: "Alice & Bob", Message: &message, IssuerName: "Inviter", App: "BiOME", Role: "contributor",
		URL: url.URL{Scheme: "https", Host: "biome.example", Path: "/register", RawQuery: "token=abc"},
	}
	var body strings.Builder
	require.NoError(t, email_templates.Invitation(data).Render(context.Background(), &body))
	require.Contains(t, body.String(), "Alice &amp; Bob")
	require.Contains(t, body.String(), "&lt;script&gt;")
	require.NotContains(t, body.String(), "<script>")
	require.Contains(t, body.String(), "Inviter")
	require.Contains(t, body.String(), "https://biome.example/register?token=abc")
}
