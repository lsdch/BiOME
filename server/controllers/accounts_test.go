package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/middleware"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
	"github.com/stretchr/testify/require"
)

type accountServiceStub struct {
	*services.AccountService
	err                               error
	userID                            uuid.UUID
	calls                             int
	identifier, path, token, password string
}

func (s *accountServiceStub) UpdatePassword(_ context.Context, _ *db.DB, id uuid.UUID, params models.PasswordUpdateParams) error {
	s.calls++
	s.userID, s.password = id, params.NewPassword
	s.identifier = params.OldPassword
	return s.err
}

func TestUpdatePasswordEndpoint(t *testing.T) {
	r := router.New(gin.New(), "", huma.DefaultConfig("test", "1"))
	service := &accountServiceStub{}
	var session *models.SessionContext
	r.API.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		if session != nil {
			ctx = auth.WithSession(ctx, session)
		}
		next(ctx)
	})
	c := &AccountsController{service: service}
	c.RegisterRoutes(&r)
	op := r.API.OpenAPI().Paths["/accounts/password"].Put
	policy, ok := middleware.GetAccessPolicy(*op)
	require.True(t, ok)
	require.True(t, policy.IsAuthRequired())
	require.False(t, policy.Eval(nil))
	request := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPut, "/accounts/password", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		r.Engine.ServeHTTP(response, req)
		return response
	}
	body := `{"old_password":"old-secret","new_password":"new-secret"}`
	require.Equal(t, http.StatusUnauthorized, request(body).Code)
	require.Zero(t, service.calls)
	session = &models.SessionContext{User: models.User{ID: uuid.New()}}
	response := request(body)
	require.Equal(t, http.StatusNoContent, response.Code)
	require.Empty(t, response.Body.String())
	require.Equal(t, session.User.ID, service.userID)
	require.Equal(t, "old-secret", service.identifier)
	require.Equal(t, "new-secret", service.password)
	for _, body := range []string{`{}`, `{"old_password":"old"}`, `{"old_password":"","new_password":"new"}`} {
		calls := service.calls
		require.Equal(t, http.StatusUnprocessableEntity, request(body).Code)
		require.Equal(t, calls, service.calls)
	}
	for _, tc := range []struct {
		err    error
		status int
	}{
		{services.ErrInvalidCurrentPassword, http.StatusForbidden},
		{services.ErrInvalidPasswordUpdate, http.StatusUnprocessableEntity},
		{fmt.Errorf("private SQL error"), http.StatusInternalServerError},
	} {
		service.err = fmt.Errorf("private details: %w", tc.err)
		response := request(body)
		require.Equal(t, tc.status, response.Code)
		require.NotContains(t, response.Body.String(), "private")
	}
}

func (s *accountServiceStub) RequestPasswordReset(_ context.Context, _ *db.DB, identifier, path string) error {
	s.calls++
	s.identifier, s.path = identifier, path
	return s.err
}

func (s *accountServiceStub) ResetPassword(_ context.Context, _ *db.DB, token string, params models.PasswordResetParams) error {
	s.calls++
	s.token, s.password = token, params.NewPassword
	return s.err
}

func (s *accountServiceStub) VerifyPasswordResetToken(_ context.Context, _ db.Querier, token string) error {
	s.calls++
	s.token = token
	return s.err
}

func TestAccountsPasswordResetEndpoints(t *testing.T) {
	r := router.New(gin.New(), "", huma.DefaultConfig("test", "1"))
	service := &accountServiceStub{}
	c := &AccountsController{service: service}
	c.RegisterRoutes(&r)
	for _, path := range []string{"/accounts/password-reset/request", "/accounts/password-reset", "/accounts/password-reset/verify"} {
		op := r.API.OpenAPI().Paths[path].Post
		require.NotNil(t, op)
		policy, ok := middleware.GetAccessPolicy(*op)
		require.True(t, ok)
		require.False(t, policy.IsAuthRequired())
		require.True(t, policy.Eval(nil))
	}
	request := func(path, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		r.Engine.ServeHTTP(response, req)
		return response
	}
	response := request("/accounts/password-reset/request", `{"email":"alice@example.org","client_path":"password-reset"}`)
	require.Equal(t, http.StatusNoContent, response.Code, response.Body.String())
	require.Empty(t, response.Body.String())
	require.Equal(t, "alice@example.org", service.identifier)
	require.Equal(t, "password-reset", service.path)
	response = request("/accounts/password-reset", `{"token":"secret-token","new_password":"secret-password"}`)
	require.Equal(t, http.StatusNoContent, response.Code, response.Body.String())
	require.Empty(t, response.Body.String())
	require.Equal(t, "secret-token", service.token)
	require.Equal(t, "secret-password", service.password)
	response = request("/accounts/password-reset/verify", `{"token":"verification-token"}`)
	require.Equal(t, http.StatusNoContent, response.Code, response.Body.String())
	require.Empty(t, response.Body.String())
	require.Equal(t, "verification-token", service.token)
	for _, tc := range []struct{ path, body string }{
		{"/accounts/password-reset/verify", `{}`},
		{"/accounts/password-reset/verify", `{"token":""}`},
		{"/accounts/password-reset/request", `{}`},
		{"/accounts/password-reset/request", `{"identifier":""}`},
		{"/accounts/password-reset", `{"new_password":"password"}`},
		{"/accounts/password-reset", `{"token":"token"}`},
		{"/accounts/password-reset", `{"token":"","new_password":""}`},
	} {
		calls := service.calls
		require.Equal(t, http.StatusUnprocessableEntity, request(tc.path, tc.body).Code)
		require.Equal(t, calls, service.calls)
	}
	for _, tc := range []struct {
		path   string
		err    error
		status int
	}{
		{"/accounts/password-reset/verify", services.ErrInvalidPasswordResetToken, http.StatusBadRequest},
		{"/accounts/password-reset/verify", fmt.Errorf("SQL secret"), http.StatusInternalServerError},
		{"/accounts/password-reset", services.ErrInvalidPasswordResetToken, http.StatusBadRequest},
		{"/accounts/password-reset", services.ErrInvalidPasswordReset, http.StatusUnprocessableEntity},
		{"/accounts/password-reset/request", services.ErrMailerUnavailable, http.StatusServiceUnavailable},
		{"/accounts/password-reset", fmt.Errorf("SQL secret"), http.StatusInternalServerError},
		{"/accounts/password-reset/request", fmt.Errorf("SMTP secret"), http.StatusInternalServerError},
	} {
		service.err = fmt.Errorf("private details: %w", tc.err)
		body := `{"token":"token","new_password":"password"}`
		if tc.path == "/accounts/password-reset/request" {
			body = `{"email":"alice@example.org","client_path":"password-reset"}`
		}
		if tc.path == "/accounts/password-reset/verify" {
			body = `{"token":"token"}`
		}
		response := request(tc.path, body)
		require.Equal(t, tc.status, response.Code, response.Body.String())
		require.NotContains(t, response.Body.String(), "private details")
		require.NotContains(t, response.Body.String(), "secret")
	}
}
