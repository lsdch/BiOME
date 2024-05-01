package settings

import (
	"context"
	"darco/proto/models/settings"
	"darco/proto/resolvers"
	"darco/proto/router"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/disintegration/imaging"
)

func RegisterRoutes(r router.Router) {
	api := r.RouteGroup("/settings").WithTags([]string{"Settings"})

	router.Register(api, "InstanceSettings",
		huma.Operation{
			Path:    "/instance",
			Method:  http.MethodGet,
			Summary: "Instance settings",
		}, GetInstanceSettings)
	router.Register(api, "UpdateInstanceSettings",
		huma.Operation{
			Path:    "/instance",
			Method:  http.MethodPost,
			Summary: "Update instance settings",
		}, UpdateInstanceSettings)

	router.Register(api, "SecuritySettings",
		huma.Operation{
			Path:    "/security",
			Method:  http.MethodGet,
			Summary: "Security settings",
		}, GetSecuritySettings)
	router.Register(api, "UpdateSecuritySettings",
		huma.Operation{
			Path:    "/security",
			Method:  http.MethodPost,
			Summary: "Update security settings",
		}, UpdateSecuritySettings)

	router.Register(api, "EmailSettings",
		huma.Operation{
			Path:    "/emailing",
			Method:  http.MethodGet,
			Summary: "Email settings",
		}, GetEmailSettings)
	router.Register(api, "UpdateEmailSettings",
		huma.Operation{
			Path:    "/emailing",
			Method:  http.MethodPost,
			Summary: "Update email settings",
		}, UpdateEmailSettings)
	router.Register(api, "TestSMTP",
		huma.Operation{
			Path:    "/emailing/test-dial",
			Method:  http.MethodPost,
			Summary: "Test SMTP connection",
		}, TestSMTP)

	router.Register(api, "SetAppIcon",
		huma.Operation{
			Path:    "/icon",
			Method:  http.MethodPost,
			Summary: "Set app icon",
		}, SetAppIcon)
}

type InstanceSettings struct{ Body settings.InstanceSettings }

func GetInstanceSettings(ctx context.Context, input *struct{}) (*InstanceSettings, error) {
	return &InstanceSettings{Body: settings.Instance()}, nil
}
func UpdateInstanceSettings(ctx context.Context,
	input *struct {
		resolvers.AccessRestricted[resolvers.Admin]
		Body settings.InstanceSettingsInput
	},
) (*InstanceSettings, error) {
	updated, err := input.Body.Save(input.DB())
	return &InstanceSettings{Body: *updated}, err
}

type SecuritySettings struct{ Body settings.SecuritySettings }

func GetSecuritySettings(ctx context.Context, input *resolvers.AccessRestricted[resolvers.Admin]) (*SecuritySettings, error) {
	return &SecuritySettings{Body: settings.Security()}, nil
}

func UpdateSecuritySettings(ctx context.Context,
	input *struct {
		resolvers.AccessRestricted[resolvers.Admin]
		Body settings.SecuritySettingsInput
	},
) (*SecuritySettings, error) {
	updated, err := input.Body.Save(input.DB())
	return &SecuritySettings{Body: *updated}, err
}

type EmailSettings struct{ Body settings.EmailSettings }

func GetEmailSettings(ctx context.Context, input *resolvers.AccessRestricted[resolvers.Admin]) (*EmailSettings, error) {
	return &EmailSettings{Body: settings.Email()}, nil
}

type EmailSettingsInput struct {
	resolvers.AccessRestricted[resolvers.Admin]
	Body settings.EmailSettingsInput
}

type SMTPConnectionStatus struct{ Body bool }

func UpdateEmailSettings(ctx context.Context, input *EmailSettingsInput) (*EmailSettings, error) {
	if status, err := TestSMTP(ctx, input); !status.Body {
		return nil, err
	}
	updated, err := input.Body.Save(input.DB())
	return &EmailSettings{Body: *updated}, err
}

func TestSMTP(ctx context.Context, input *EmailSettingsInput) (*SMTPConnectionStatus, error) {
	if err := input.Body.TestConnection(); err != nil {
		return &SMTPConnectionStatus{false}, huma.Error422UnprocessableEntity("SMTP connection failed", err)
	}
	return &SMTPConnectionStatus{true}, nil
}

// router.Register(api, "SetAppIcon",
// 	huma.Operation{
// 		Path:          "/icon",
// 		Method:        http.MethodPost,
// 		Summary:       "Set app icon",
// 		DefaultStatus: http.StatusSeeOther,
// 		// RequestBody: &huma.RequestBody{
// 		// 	Required:    true,
// 		// 	Description: "Icon file encoded as multipart form data",
// 		// 	Content: map[string]*huma.MediaType{
// 		// 		"multipart/form-data": {
// 		// 			Schema: &huma.Schema{
// 		// 				Type: "object",
// 		// 				Properties: map[string]*huma.Schema{
// 		// 					FORMDATA_ICON_KEY: {
// 		// 						Type:        "string",
// 		// 						Format:      "binary",
// 		// 						Description: "Icon file to upload",
// 		// 					},
// 		// 				},
// 		// 			},
// 		// 		},
// 		// 	},
// 		// },
// 	}, SetAppIcon)

// type AppIconInputData struct {
// 	Icon multipart.File `form-data:"icon" content-type:"image/png" required:"true"`
// 	// Test []multipart.File `form-data:"test" content-type:"image/png" required:"true"`
// }

type AppIconInput struct {
	// resolvers.AccessRestricted[resolvers.Admin]
	resolvers.HostResolver
	RawBody multipart.Form
	// RawBody huma.MultipartFormFiles[AppIconInputData]
	Image image.Image
}

const (
	APP_ICON_PATH     = "assets/app_icon.png"
	FORMDATA_ICON_KEY = "iconFile"
)

func (i *AppIconInput) Resolve(ctx huma.Context) []error {
	fileHeaders := i.RawBody.File[FORMDATA_ICON_KEY]
	if len(fileHeaders) == 0 {
		return []error{&huma.ErrorDetail{Message: "Empty file content received"}}
	} else if len(fileHeaders) > 1 {
		return []error{&huma.ErrorDetail{Message: "Multiple files received"}}
	}
	f, err := fileHeaders[0].Open()
	if err != nil {
		return []error{&huma.ErrorDetail{Message: "Failed to open submitted file"}}
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return []error{&huma.ErrorDetail{Message: "Failed to decode image from file. Accepted formats are: PNG, JPEG."}}
	}
	i.Image = img
	return nil
}

// func (i *AppIconInput) Resolve(ctx huma.Context) []error {
// 	formData := i.RawBody.Data()

// 	img, _, err := image.Decode(formData.Icon)
// 	if err != nil {
// 		return []error{&huma.ErrorDetail{Message: "Failed to decode image from file. Accepted formats are: PNG, JPEG."}}
// 	}
// 	i.Image = img
// 	return nil
// }

type SetAppIconOutput struct {
	Location url.URL `header:"Location" format:"uri"`
}

func SetAppIcon(ctx context.Context, input *AppIconInput) (*SetAppIconOutput, error) {

	resizedImg := imaging.Resize(input.Image, 300, 300, imaging.Lanczos)

	writer, err := os.Create("assets/app_icon.png")
	if err != nil {
		return nil, err
	}
	defer writer.Close()
	if err := png.Encode(writer, resizedImg); err != nil {
		return nil, err
	}

	return &SetAppIconOutput{
		Location: input.GenerateURL("assets/app_icon.png"),
	}, nil
}
