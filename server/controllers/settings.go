package controllers

import (
	"context"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"net/http"
	"net/url"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/disintegration/imaging"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
)

type SettingsController struct {
	db      *db.DB
	service *services.SettingsService
}

func NewSettingsController(db *db.DB, service *services.SettingsService) *SettingsController {
	return &SettingsController{db: db, service: service}
}

func (c *SettingsController) GetInstanceSettings(
	ctx context.Context,
	input *struct{},
) (*BodyTransporter[models.InstanceSettings], error) {
	settings := c.service.GetSettings()
	return &BodyTransporter[models.InstanceSettings]{Body: settings}, nil
}

func (c *SettingsController) UpdateInstanceSettings(
	ctx context.Context,
	input *BodyTransporter[models.InstanceSettingsUpdate],
) (*struct{}, error) {
	err := c.service.UpdateInstanceSettings(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *SettingsController) TestSMTP(
	ctx context.Context,
	input *struct{},
) (*BodyTransporter[bool], error) {
	status, err := c.service.TestSMTPConnection(ctx)
	if err != nil {
		return nil, err
	}
	return &BodyTransporter[bool]{Body: status}, nil
}

func (c *SettingsController) TogglePublicAccess(
	ctx context.Context,
	input *BodyTransporter[bool],
) (*struct{}, error) {
	err := c.service.TogglePublicAccess(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *SettingsController) TogglePublicRegistration(
	ctx context.Context,
	input *BodyTransporter[bool],
) (*struct{}, error) {
	err := c.service.TogglePublicRegistration(ctx, c.db, input.Body)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type AppIconInputData struct {
	Icon huma.FormFile `form:"icon" contentType:"image/png,image/jpeg" required:"true"`
}

type AppIconInput struct {
	RawBody     huma.MultipartFormFiles[AppIconInputData]
	Image       image.Image
	ContentType string // Uploaded image format
}

func (i *AppIconInput) Resolve(ctx huma.Context) []error {
	formData := i.RawBody.Data()

	img, _, err := image.Decode(formData.Icon)
	if err != nil {
		return []error{&huma.ErrorDetail{Message: "Failed to decode image from file. Accepted formats are: PNG, JPEG."}}
	}
	i.Image = img
	return nil
}

type SetAppIconOutput struct {
	Location url.URL `header:"Location" format:"uri"`
}

func (c *SettingsController) SetAppIcon(ctx context.Context, input *AppIconInput) (*SetAppIconOutput, error) {

	resizedImg := imaging.Resize(input.Image, 300, 300, imaging.Lanczos)

	writer, err := os.Create("assets/app_icon.png")
	if err != nil {
		return nil, err
	}
	defer writer.Close()
	err = png.Encode(writer, resizedImg)
	if err != nil {
		return nil, err
	}

	return &SetAppIconOutput{
		Location: *c.service.Config.AppPublicBaseURL.JoinPath("assets/app_icon.png"),
	}, nil
}

func (c *SettingsController) RegisterRoutes(r *router.Router) {
	settingsAPI := r.RouteGroup("/settings").WithTags([]string{"Settings"})

	router.NewSpec(
		settingsAPI,
		"GetInstanceSettings",
		huma.Operation{
			Path:    "/instance",
			Method:  http.MethodGet,
			Summary: "Get instance settings",
		},
		c.GetInstanceSettings,
	).
		WithAccessPolicy(auth.Public()).
		Register(r)

	router.NewSpec(
		settingsAPI,
		"UpdateInstanceSettings",
		huma.Operation{
			Path:    "/instance",
			Method:  http.MethodPut,
			Summary: "Update instance settings",
		},
		c.UpdateInstanceSettings,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).
		Register(r)

	router.NewSpec(
		settingsAPI,
		"TestSMTPConnection",
		huma.Operation{
			Path:    "/smtp/test",
			Method:  http.MethodGet,
			Summary: "Test SMTP connection",
		},
		c.TestSMTP,
	).
		WithAccessPolicy(auth.Authenticated()).
		Register(r)

	router.NewSpec(
		settingsAPI,
		"TogglePublicAccess",
		huma.Operation{
			Path:    "/instance/public-access",
			Method:  http.MethodPut,
			Summary: "Toggle public access to the instance",
		},
		c.TogglePublicAccess,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).
		Register(r)

	router.NewSpec(
		settingsAPI,
		"TogglePublicRegistration",
		huma.Operation{
			Path:    "/instance/public-registration",
			Method:  http.MethodPut,
			Summary: "Toggle public registration for the instance",
		},
		c.TogglePublicRegistration,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).
		Register(r)

	router.NewSpec(
		settingsAPI,
		"SetAppIcon",
		huma.Operation{
			Path:    "/instance/app-icon",
			Method:  http.MethodPost,
			Summary: "Set the application icon",
		},
		c.SetAppIcon,
	).
		WithAccessPolicy(auth.Role(biomedb.UserRoleAdmin)).
		Register(r)
}
