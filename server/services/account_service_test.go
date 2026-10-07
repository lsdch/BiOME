package services

import (
	"context"
	"errors"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type resetMailer struct {
	to, body string
	err      error
}

func (*resetMailer) Available() bool { return true }
func (m *resetMailer) Send(ctx context.Context, to, subject string, template templ.Component) error {
	m.to = to
	var body strings.Builder
	if err := template.Render(ctx, &body); err != nil {
		return err
	}
	m.body = body.String()
	return m.err
}

func TestPasswordResetEarlyFailures(t *testing.T) {
	s := &AccountService{}
	require.ErrorIs(t, s.UpdatePassword(context.Background(), nil, uuid.Nil, models.PasswordUpdateParams{}), ErrInvalidCurrentPassword)
	require.ErrorIs(t, s.VerifyPasswordResetToken(context.Background(), nil, ""), ErrInvalidPasswordResetToken)
	require.ErrorIs(t, s.ResetPassword(context.Background(), nil, "", models.PasswordResetParams{}), ErrInvalidPasswordResetToken)
	require.ErrorIs(t, s.RequestPasswordReset(context.Background(), nil, "alice", "reset"), ErrMailerUnavailable)
	s.emailService = &resetMailer{}
	require.ErrorContains(t, s.RequestPasswordReset(context.Background(), nil, "alice", "reset"), "URL configuration")
	s.publicBaseURL = url.URL{Scheme: "https", Host: "biome.example"}
	require.ErrorContains(t, s.RequestPasswordReset(context.Background(), nil, "alice", " "), "URL configuration")
}

func TestUpdatePassword(t *testing.T) {
	database, pool := passwordResetTestDB(t)
	ctx := context.Background()
	s := &AccountService{}
	old := "Old!Cactus-87.River$Quartz"
	password := "New!Cactus-92.River$Quartz"
	hash, err := s.HashPassword(old)
	require.NoError(t, err)
	u, err := database.Queries().CreateUser(ctx, biomedb.CreateUserParams{
		Login: "alice", Email: "alice@example.org", FirstName: "Alice", LastName: "Test", Role: "visitor", PasswordHash: string(hash),
	})
	require.NoError(t, err)
	require.NoError(t, database.Queries().CreatePasswordResetToken(ctx, biomedb.CreatePasswordResetTokenParams{UserID: u.ID, TokenHash: HashToken("reset"), ExpiresAt: time.Now().Add(time.Hour)}))
	_, err = pool.Exec(ctx, "INSERT INTO user_sessions (id, user_id, refresh_token_hash, expires_at) VALUES ($1, $2, 'refresh', now() + interval '1 hour')", uuid.New(), u.ID)
	require.NoError(t, err)
	require.ErrorIs(t, s.UpdatePassword(ctx, database, u.ID, models.PasswordUpdateParams{OldPassword: "wrong", NewPassword: password}), ErrInvalidCurrentPassword)
	for _, invalid := range []string{"", "weak", strings.Repeat("é", 37)} {
		require.ErrorIs(t, s.UpdatePassword(ctx, database, u.ID, models.PasswordUpdateParams{OldPassword: old, NewPassword: invalid}), ErrInvalidPasswordUpdate)
	}
	params := models.PasswordUpdateParams{OldPassword: old, NewPassword: password}
	_, err = pool.Exec(ctx, "ALTER TABLE user_sessions ADD CONSTRAINT reject_update_revocation CHECK (revoked_at IS NULL)")
	require.NoError(t, err)
	require.Error(t, s.UpdatePassword(ctx, database, u.ID, params))
	stored, err := database.Queries().GetUserByID(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, string(hash), stored.PasswordHash)
	require.NoError(t, s.VerifyPasswordResetToken(ctx, database, "reset"))
	_, err = pool.Exec(ctx, "ALTER TABLE user_sessions DROP CONSTRAINT reject_update_revocation")
	require.NoError(t, err)
	require.NoError(t, s.UpdatePassword(ctx, database, u.ID, params))
	stored, err = database.Queries().GetUserByID(ctx, u.ID)
	require.NoError(t, err)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte(password)))
	require.ErrorIs(t, s.UpdatePassword(ctx, database, u.ID, params), ErrInvalidCurrentPassword)
	require.ErrorIs(t, s.VerifyPasswordResetToken(ctx, database, "reset"), ErrInvalidPasswordResetToken)
	var active int
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM user_sessions WHERE revoked_at IS NULL").Scan(&active))
	require.Zero(t, active)
}

func passwordResetTestDB(t *testing.T) (*db.DB, *pgxpool.Pool) {
	t.Helper()
	dsn := os.Getenv("BIOME_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("BIOME_TEST_DATABASE_URL is required")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(admin.Close)
	schema := "password_reset_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	_, err = admin.Exec(ctx, "CREATE SCHEMA "+schema)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := admin.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
		require.NoError(t, err)
	})
	cfg, err := pgxpool.ParseConfig(dsn)
	require.NoError(t, err)
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err)
	database := db.NewDB(pool)
	t.Cleanup(database.Close)
	_, err = pool.Exec(ctx, "CREATE TYPE user_role AS ENUM ('visitor')")
	require.NoError(t, err)
	users, err := os.ReadFile("../db/schema/010_users.sql")
	require.NoError(t, err)
	// TEXT avoids installing citext in the shared test database.
	_, err = pool.Exec(ctx, strings.ReplaceAll(string(users), "CITEXT", "TEXT"))
	require.NoError(t, err)
	migration, err := os.ReadFile("../db/migrations/001_password_reset_tokens.up.sql")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, string(migration))
	require.NoError(t, err)
	sessions, err := os.ReadFile("../db/schema/013_user_sessions.sql")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, string(sessions))
	require.NoError(t, err)
	return database, pool
}

func TestPasswordResetRequest(t *testing.T) {
	database, pool := passwordResetTestDB(t)
	ctx := context.Background()
	user, err := database.Queries().CreateUser(ctx, biomedb.CreateUserParams{
		Login: "alice", Email: "alice@example.org", FirstName: "Alice", LastName: "Test", PasswordHash: "unchanged", Role: "visitor",
	})
	require.NoError(t, err)
	mailer := &resetMailer{}
	s := &AccountService{emailService: mailer, publicBaseURL: url.URL{Scheme: "https", Host: "biome.example", Path: "/app"}}
	require.NoError(t, s.RequestPasswordReset(ctx, database, "unknown", "reset"))
	require.Empty(t, mailer.to)
	require.NoError(t, s.RequestPasswordReset(ctx, database, user.Login, "reset"))
	require.Equal(t, user.Email, mailer.to)
	token := regexp.MustCompile(`token=([A-Za-z0-9_-]+)`).FindStringSubmatch(mailer.body)
	require.Len(t, token, 2)
	require.Contains(t, mailer.body, "https://biome.example/app/reset?token=")
	var hash string
	var created, expires time.Time
	err = pool.QueryRow(ctx, "SELECT token_hash, created_at, expires_at FROM password_reset_tokens WHERE user_id = $1", user.ID).Scan(&hash, &created, &expires)
	require.NoError(t, err)
	require.Equal(t, HashToken(token[1]), hash)
	require.WithinDuration(t, created.Add(time.Hour), expires, 5*time.Second)
	stored, err := database.Queries().GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "unchanged", stored.PasswordHash)
	mailer.err = errors.New("delivery failed")
	require.ErrorIs(t, s.RequestPasswordReset(ctx, database, user.Email, "reset"), mailer.err)
	_, err = pool.Exec(ctx, "UPDATE users SET active = false WHERE id = $1", user.ID)
	require.NoError(t, err)
	mailer.to = ""
	require.NoError(t, s.RequestPasswordReset(ctx, database, user.Login, "reset"))
	require.Empty(t, mailer.to)
	var count int
	err = pool.QueryRow(ctx, "SELECT count(*) FROM password_reset_tokens").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestPasswordResetConsumption(t *testing.T) {
	database, pool := passwordResetTestDB(t)
	ctx := context.Background()
	s := &AccountService{}
	params := models.PasswordResetParams{NewPassword: "Lunar!Cactus-87.River$Quartz"}
	user, err := database.Queries().CreateUser(ctx, biomedb.CreateUserParams{
		Login: "alice", Email: "alice@example.org", FirstName: "Alice", LastName: "Test", PasswordHash: "unchanged", Role: "visitor",
	})
	require.NoError(t, err)
	newToken := func() string {
		token, err := GenerateSecureToken()
		require.NoError(t, err)
		require.NoError(t, database.Queries().CreatePasswordResetToken(ctx, biomedb.CreatePasswordResetTokenParams{
			UserID: user.ID, TokenHash: HashToken(token), ExpiresAt: time.Now().Add(time.Hour),
		}))
		return token
	}
	token, other := newToken(), newToken()
	// Repeated verification must preserve the entire token row and allow reset.
	var before, after string
	require.NoError(t, pool.QueryRow(ctx, "SELECT row_to_json(t)::text FROM password_reset_tokens t WHERE token_hash = $1", HashToken(token)).Scan(&before))
	for range 2 {
		require.NoError(t, s.VerifyPasswordResetToken(ctx, database, token))
	}
	require.NoError(t, pool.QueryRow(ctx, "SELECT row_to_json(t)::text FROM password_reset_tokens t WHERE token_hash = $1", HashToken(token)).Scan(&after))
	require.Equal(t, before, after)
	require.ErrorIs(t, s.VerifyPasswordResetToken(ctx, database, "unknown"), ErrInvalidPasswordResetToken)
	require.ErrorIs(t, s.ResetPassword(ctx, database, "unknown", params), ErrInvalidPasswordResetToken)
	require.ErrorIs(t, s.ResetPassword(ctx, database, token, models.PasswordResetParams{NewPassword: "weak"}), ErrInvalidPasswordReset)
	_, err = pool.Exec(ctx, "INSERT INTO user_sessions (id, user_id, refresh_token_hash, expires_at) VALUES ($1, $2, 'refresh', now() + interval '1 hour')", uuid.New(), user.ID)
	require.NoError(t, err)
	// A late failure must roll back both password mutation and token consumption.
	_, err = pool.Exec(ctx, "ALTER TABLE user_sessions ADD CONSTRAINT reject_revocation CHECK (revoked_at IS NULL)")
	require.NoError(t, err)
	require.Error(t, s.ResetPassword(ctx, database, token, params))
	stored, err := database.Queries().GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "unchanged", stored.PasswordHash)
	_, err = pool.Exec(ctx, "ALTER TABLE user_sessions DROP CONSTRAINT reject_revocation")
	require.NoError(t, err)
	require.NoError(t, s.ResetPassword(ctx, database, token, params))
	stored, err = database.Queries().GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte(params.NewPassword)))
	require.ErrorIs(t, s.ResetPassword(ctx, database, token, params), ErrInvalidPasswordResetToken)
	require.ErrorIs(t, s.ResetPassword(ctx, database, other, params), ErrInvalidPasswordResetToken)
	require.ErrorIs(t, s.VerifyPasswordResetToken(ctx, database, token), ErrInvalidPasswordResetToken)
	require.ErrorIs(t, s.VerifyPasswordResetToken(ctx, database, other), ErrInvalidPasswordResetToken)
	var activeSessions int
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM user_sessions WHERE revoked_at IS NULL").Scan(&activeSessions))
	require.Zero(t, activeSessions)
	expired := newToken()
	_, err = pool.Exec(ctx, "UPDATE password_reset_tokens SET created_at = now() - interval '2 hours', expires_at = now() - interval '1 hour' WHERE token_hash = $1", HashToken(expired))
	require.NoError(t, err)
	require.ErrorIs(t, s.ResetPassword(ctx, database, expired, params), ErrInvalidPasswordResetToken)
	require.ErrorIs(t, s.VerifyPasswordResetToken(ctx, database, expired), ErrInvalidPasswordResetToken)
	inactive := newToken()
	_, err = pool.Exec(ctx, "UPDATE users SET active = false WHERE id = $1", user.ID)
	require.NoError(t, err)
	require.ErrorIs(t, s.ResetPassword(ctx, database, inactive, params), ErrInvalidPasswordResetToken)
	require.ErrorIs(t, s.VerifyPasswordResetToken(ctx, database, inactive), ErrInvalidPasswordResetToken)
	_, err = pool.Exec(ctx, "UPDATE users SET active = true WHERE id = $1", user.ID)
	require.NoError(t, err)
	for _, differentTokens := range []bool{false, true} {
		first, second := newToken(), newToken()
		if !differentTokens {
			second = first
		}
		start := make(chan struct{})
		results := make(chan error, 2)
		var workers sync.WaitGroup
		for _, candidate := range []string{first, second} {
			workers.Add(1)
			go func(token string) {
				defer workers.Done()
				<-start
				results <- s.ResetPassword(ctx, database, token, params)
			}(candidate)
		}
		close(start)
		workers.Wait()
		close(results)
		successes := 0
		for err := range results {
			if err == nil {
				successes++
			} else {
				require.ErrorIs(t, err, ErrInvalidPasswordResetToken)
			}
		}
		require.Equal(t, 1, successes)
	}
}
