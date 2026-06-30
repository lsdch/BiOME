package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/url"

	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
)

type RegistrationService struct {
	db           *db.DB
	config       config.Config
	emailService *EmailService
}

func NewRegistrationService(db *db.DB, config config.Config, emailService *EmailService) *RegistrationService {
	return &RegistrationService{db: db, config: config, emailService: emailService}
}

func (s *RegistrationService) CreateInvitation(
	ctx context.Context,
	params models.CreateInvitationParams,
) (models.InvitationResult, error) {

	var result models.InvitationResult

	err := s.db.WithTx(ctx, func(q *biomedb.Queries) error {

		invitation, err := q.CreateInvitation(ctx, biomedb.CreateInvitationParams{
			Email:       params.Email,
			InviteeName: params.InviteeName,
			Role:        params.Role,
			Message:     params.Message.ToPtr(),
			InviterID:   models.UUIDToPg(params.InviterID),
			ExpiresAt:   params.ExpiresAt,
		})
		if err != nil {
			return err
		}

		rawToken, err := GenerateSecureToken()
		if err != nil {
			return err
		}

		hash := HashToken(rawToken)

		_, err = q.CreateInvitationToken(ctx, invitation.ID, hash)
		if err != nil {
			return err
		}

		result = models.InvitationResult{
			Invitation: invitation,
			Token:      rawToken,
		}

		return nil
	})

	return result, err
}

func (s *RegistrationService) InvitationURL(clientPath string, token string) *url.URL {
	u := s.config.AppPublicBaseURL.JoinPath(clientPath)
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u
}

// func (s *RegistrationService) SendInvitationEmail(
// 	ctx context.Context,
// 	invitation InvitationResult,
// ) error {

// 	templateData := email_templates.InvitationData{
// 		Name:       invitation.Invitation.InviteeName,
// 		IssuerName: invitation.Invitation.InviterName,
// 		App:        settings.Instance().Name,
// 		Role:       string(invitation.Invitation.Role),
// 		URL:        *s.InvitationURL("register", invitation.Token),
// 	}

// 	s.emailService.Send(ctx,
// 		invitation.Invitation.Email,
// 		s.emailService.,
// 		"You're invited to join Biome",
// 		email_templates.Invitation(templateData),
// 	)
// 	return nil
// }

func GenerateSecureToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *RegistrationService) GetInvitationByToken(
	ctx context.Context,
	token string,
) (biomedb.Invitation, error) {

	return s.db.Queries().GetInvitationByTokenHash(
		ctx,
		HashToken(token),
	)
}

func (s *RegistrationService) RegisterFromInvitation(
	ctx context.Context,
	token string,
	params models.RegisterFromInvitationParams,
) (biomedb.User, error) {

	tokenHash := HashToken(token)
	u, err := s.db.Queries().CreateUserFromInvitationToken(
		ctx,
		params.ToParams(tokenHash),
	)
	if err != nil {
		return biomedb.User{}, err
	}
	return u, nil
}
