package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	email_templates "github.com/lsdch/biome/templates"
)

var (
	ErrInvalidInvitation       = errors.New("invitation is invalid, expired, cancelled or already redeemed")
	ErrInvalidRegistration     = errors.New("invalid registration parameters")
	ErrAccountAlreadyExists    = errors.New("an account already exists for this email or login")
	ErrInvitationAlreadyExists = errors.New("a pending invitation already exists for this email")
)

type RegistrationService struct {
	config       config.Config
	emailService *EmailService
}

func NewRegistrationService(config config.Config, emailService *EmailService) *RegistrationService {
	return &RegistrationService{config: config, emailService: emailService}
}

// CreateInvitation persists the invitation and its hashed token atomically.
// When q is a transaction, its caller owns the commit. Send email only after commit.
func (s *RegistrationService) CreateInvitation(ctx context.Context, q db.Querier, params models.CreateInvitationParams) (models.InvitationResult, error) {
	address, err := mail.ParseAddress(params.Email)
	if err != nil || address.Address != params.Email || strings.TrimSpace(params.InviteeName) == "" || !params.Role.Valid() || !params.ExpiresAt.After(time.Now()) {
		return models.InvitationResult{}, fmt.Errorf("%w: invalid invitation email, name, role or expiry", ErrInvalidRegistration)
	}
	rawToken, err := GenerateSecureToken()
	if err != nil {
		return models.InvitationResult{}, err
	}
	var result models.InvitationResult
	err = q.WithTx(ctx, func(tx *db.Tx) error {
		exists, err := tx.Queries().UserExistsByEmail(ctx, params.Email)
		if err != nil {
			return err
		}
		if exists {
			return ErrAccountAlreadyExists
		}
		// Expiration does not automatically update the partial unique index.
		if err := tx.Queries().ExpireInvitationsForEmail(ctx, params.Email); err != nil {
			return err
		}
		invitation, err := tx.Queries().CreateInvitation(ctx, biomedb.CreateInvitationParams{
			Email: params.Email, InviteeName: params.InviteeName, Role: params.Role,
			Message: params.Message.ToPtr(), InviterID: models.UUIDToPg(params.InviterID), ExpiresAt: params.ExpiresAt,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "invitations_one_active_per_email_idx" {
				return fmt.Errorf("%w: %w", ErrInvitationAlreadyExists, err)
			}
			return err
		}
		if _, err := tx.Queries().CreateInvitationToken(ctx, invitation.ID, HashToken(rawToken)); err != nil {
			return err
		}
		result = models.InvitationResult{Invitation: invitation, Token: rawToken}
		return nil
	})
	if err != nil {
		return models.InvitationResult{}, fmt.Errorf("create invitation: %w", err)
	}
	return result, nil
}

func (s *RegistrationService) InvitationURL(clientPath string, token string) *url.URL {
	u := s.config.AppPublicBaseURL.JoinPath(clientPath)
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u
}

// SendInvitationEmail can be retried with the same result after a delivery failure.
// It reads authoritative invitation and instance data; the caller supplies the client route.
func (s *RegistrationService) SendInvitationEmail(ctx context.Context, q db.Querier, result models.InvitationResult, clientPath string) error {
	if s.emailService == nil || !s.emailService.Available() {
		return ErrMailerUnavailable
	}
	base := s.config.AppPublicBaseURL
	if (base.Scheme != "https" && base.Scheme != "http") || base.Host == "" || strings.TrimSpace(clientPath) == "" {
		return fmt.Errorf("invalid invitation URL configuration")
	}
	invitation, err := s.GetInvitationByToken(q, ctx, result.Token)
	if err != nil {
		return err
	}
	if invitation.ID != result.Invitation.ID {
		return ErrInvalidInvitation
	}
	settings, err := q.Queries().GetSettings(ctx)
	if err != nil {
		return fmt.Errorf("get invitation email settings: %w", err)
	}
	issuerName := settings.AppName
	if invitation.InviterID.Valid {
		inviter, err := q.Queries().GetUserByID(ctx, uuid.UUID(invitation.InviterID.Bytes))
		if err != nil {
			return fmt.Errorf("get inviter: %w", err)
		}
		issuerName = inviter.FullName
	}
	data := email_templates.InvitationData{
		IssuerName: issuerName, App: settings.AppName, Role: string(invitation.Role),
		URL: *s.InvitationURL(clientPath, result.Token), Name: invitation.InviteeName, Message: invitation.Message,
	}
	if err := s.emailService.Send(ctx, invitation.Email, data.Subject(), email_templates.Invitation(data)); err != nil {
		return fmt.Errorf("send invitation email: %w", err)
	}
	return nil
}

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

func (s *RegistrationService) GetInvitationByToken(q db.Querier, ctx context.Context, token string) (biomedb.Invitation, error) {
	if token == "" {
		return biomedb.Invitation{}, ErrInvalidInvitation
	}
	invitation, err := q.Queries().GetInvitationByTokenHash(ctx, HashToken(token))
	if errors.Is(err, pgx.ErrNoRows) {
		return biomedb.Invitation{}, ErrInvalidInvitation
	}
	return invitation, err
}

func (s *RegistrationService) CancelInvitation(ctx context.Context, q db.Querier, invitationID, revokedBy uuid.UUID) error {
	_, err := q.Queries().CancelInvitation(ctx, models.UUIDToPg(revokedBy), invitationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidInvitation
	}
	return err
}

func (s *RegistrationService) RegisterFromInvitation(q db.Querier, ctx context.Context, token string, params models.RegisterFromInvitationParams) (biomedb.User, error) {
	invitation, err := s.GetInvitationByToken(q, ctx, token)
	if err != nil {
		return biomedb.User{}, err
	}
	if err := params.Validate(invitation.Email, MIN_PASSWORD_STRENGTH); err != nil {
		return biomedb.User{}, fmt.Errorf("%w: %w", ErrInvalidRegistration, err)
	}
	input, err := params.ToParams(HashToken(token))
	if err != nil {
		return biomedb.User{}, fmt.Errorf("%w: %w", ErrInvalidRegistration, err)
	}
	// One statement locks invitation and token, inserts the account and consumes the invitation.
	u, err := q.Queries().CreateUserFromInvitationToken(ctx, input)
	if errors.Is(err, pgx.ErrNoRows) {
		return biomedb.User{}, ErrInvalidInvitation
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && (pgErr.ConstraintName == "users_email_unique" || pgErr.ConstraintName == "users_login_key") {
		return biomedb.User{}, fmt.Errorf("%w: %w", ErrAccountAlreadyExists, err)
	}
	return biomedb.User(u), err
}
