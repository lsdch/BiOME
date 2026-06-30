package accounts

import (
	"context"
	"net/http"
	"net/netip"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models/users"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"

	"github.com/danielgtaylor/huma/v2"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) GetCurrentUser(
	ctx context.Context,
	input *struct{},
) (*controllers.BodyTransporter[users.User], error) {
	session, ok := auth.SessionFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("User not found")
	}
	return &controllers.BodyTransporter[users.User]{Body: session.User}, nil
}

type LoginOutput struct {
	SessionCookie http.Cookie `header:"Set-Cookie" doc:"Session cookie storing JWT"`
	Body          services.LoginResult
}

func (c *AuthController) Login(
	ctx context.Context,
	input *struct {
		services.UserCredentials
		UserAgent string      `header:"User-Agent" doc:"User agent string of the client"`
		ClientIP  *netip.Addr `header:"client-ip" format:"ip" doc:"IP address of the client"`
	},
) (*LoginOutput, error) {
	res, err := c.authService.Login(ctx, input.UserCredentials, services.AuthOrigin{
		UserAgent: input.UserAgent,
		IP:        input.ClientIP,
	})
	if err != nil {
		return nil, err
	}
	return &LoginOutput{
		SessionCookie: c.authService.SessionCookie(res.AuthToken),
		Body:          res,
	}, nil
}

type LogoutOutput struct {
	SessionCookie http.Cookie `header:"Set-Cookie" doc:"Session cookie storing JWT"`
}

func (c *AuthController) Logout(
	ctx context.Context,
	input *struct{},
) (*LogoutOutput, error) {
	session, ok := auth.SessionFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("User not found")
	}
	err := c.authService.Logout(ctx, session.SessionID)
	if err != nil {
		return nil, err
	}
	return &LogoutOutput{
		SessionCookie: c.authService.DropSessionCookie(),
	}, nil
}

type RefreshSessionOutput struct {
	SessionCookie http.Cookie `header:"Set-Cookie" doc:"Session cookie storing JWT"`
	Body          services.SessionTokens
}

func (c *AuthController) RefreshSession(
	ctx context.Context,
	input *struct {
		RefreshToken string `header:"X-Refresh-Token" doc:"Refresh token for the session"`
	},
) (*RefreshSessionOutput, error) {
	session, ok := auth.SessionFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("User not found")
	}
	newTokens, err := c.authService.RefreshSession(ctx, session.SessionID, input.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &RefreshSessionOutput{
		SessionCookie: c.authService.SessionCookie(newTokens.AuthToken),
		Body:          newTokens,
	}, nil
}

func (c *AuthController) RegisterRoutes(r *router.Router) {
	auth_api := func(r *router.Router) router.Group {
		return r.RouteGroup("/").
			WithTags([]string{"Authentication"})
	}

	router.NewSpec(
		auth_api,
		"GetCurrentUser",
		huma.Operation{
			Method:  "GET",
			Path:    "/current-user",
			Summary: "Get the currently authenticated user",
		},
		c.GetCurrentUser,
	).WithAccessPolicy(auth.Authenticated()).Register(r)

	router.NewSpec(
		auth_api,
		"Login",
		huma.Operation{
			Method:  "POST",
			Path:    "/login",
			Summary: "Login to the application",
		},
		c.Login,
	).WithAccessPolicy(auth.Public()).Register(r)

	router.NewSpec(
		auth_api,
		"Logout",
		huma.Operation{
			Method:  "POST",
			Path:    "/logout",
			Summary: "Logout of the application",
		},
		c.Logout,
	).WithAccessPolicy(auth.Authenticated()).Register(r)

	router.NewSpec(
		auth_api,
		"RefreshSession",
		huma.Operation{
			Method:  "POST",
			Path:    "/session/refresh",
			Summary: "Refresh the authentication session",
		},
		c.RefreshSession,
	).WithAccessPolicy(auth.Authenticated()).Register(r)
}
