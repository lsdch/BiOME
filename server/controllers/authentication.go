package controllers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
)

type AuthController struct {
	db      *db.DB
	service *services.AuthService
}

func NewAuthController(db *db.DB, service *services.AuthService) *AuthController {
	return &AuthController{
		db:      db,
		service: service,
	}
}

type LoginInput struct {
	Body   models.UserCredentials
	Origin models.AuthOrigin
}

type LoginOutput struct {
	Body      models.LoginResult
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func (c *AuthController) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	var response models.LoginResult
	err := c.db.WithTx(ctx, func(tx *db.Tx) error {
		var txErr error
		response, txErr = c.service.Login(ctx, tx, input.Body, input.Origin)
		return txErr
	})
	if err != nil {
		return nil, err
	}
	return &LoginOutput{
		Body:      response,
		SetCookie: c.service.SessionCookie(response.Session.AuthToken),
	}, nil
}

type LogoutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func (c *AuthController) Logout(ctx context.Context, input *struct{}) (*LogoutOutput, error) {
	if sessionCtx, ok := auth.SessionFromContext(ctx); ok {
		err := c.service.Logout(ctx, c.db, sessionCtx.SessionID)
		if err != nil {
			return nil, err
		}
	}
	return &LogoutOutput{
		SetCookie: c.service.DropSessionCookie(),
	}, nil
}

func (c *AuthController) GetCurrentUser(ctx context.Context, input *struct{}) (*BodyTransporter[models.User], error) {
	if sessionCtx, ok := auth.SessionFromContext(ctx); ok {
		return &BodyTransporter[models.User]{Body: sessionCtx.User}, nil
	}
	return nil, nil
}

type RefreshSessionInput struct {
	RefreshToken string    `header:"X-Refresh-Token" required:"true"`
	SessionID    uuid.UUID `header:"X-Session-ID" required:"true"`
}

type RefreshSessionOutput struct {
	Body      models.SessionTokens
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func (c *AuthController) RefreshSession(ctx context.Context, input *RefreshSessionInput) (*RefreshSessionOutput, error) {
	var response models.SessionTokens
	err := c.db.WithTx(ctx, func(tx *db.Tx) error {
		var txErr error
		response, txErr = c.service.RefreshSession(ctx, tx, input.SessionID, input.RefreshToken)
		return txErr
	})
	if err != nil {
		return nil, err
	}
	return &RefreshSessionOutput{
		Body:      response,
		SetCookie: c.service.SessionCookie(response.AuthToken),
	}, nil
}

func (c *AuthController) RegisterRoutes(r *router.Router) {
	authAPI := r.RouteGroup("/auth").WithTags([]string{"Authentication"})

	router.NewSpec(
		authAPI,
		"Login",
		huma.Operation{
			Method:  "POST",
			Path:    "/login",
			Summary: "Authenticate a user and return a JWT token",
		},
		c.Login,
	).Register(r)

	router.NewSpec(
		authAPI,
		"Logout",
		huma.Operation{
			Method:  "POST",
			Path:    "/logout",
			Summary: "Invalidate the current user's session",
		},
		c.Logout,
	).WithAccessPolicy(auth.Authenticated()).Register(r)

	router.NewSpec(
		authAPI,
		"GetCurrentUser",
		huma.Operation{
			Method:  "GET",
			Path:    "/me",
			Summary: "Get the current authenticated user's information",
		},
		c.GetCurrentUser,
	).WithAccessPolicy(auth.Authenticated()).Register(r)

	router.NewSpec(
		authAPI,
		"RefreshSession",
		huma.Operation{
			Method:  "POST",
			Path:    "/refresh",
			Summary: "Refresh the current user's session",
		},
		c.RefreshSession,
	).WithAccessPolicy(auth.Public()).Register(r)
}
