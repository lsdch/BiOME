package settings

import (
	"context"
	"image"
	"image/png"
	"net/url"
	"os"

	"github.com/lsdch/biome/resolvers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/disintegration/imaging"
)

type AppIconInputData struct {
	Icon huma.FormFile `form:"icon" contentType:"image/png,image/jpeg" required:"true"`
}

type AppIconInput struct {
	resolvers.AccessRestricted[resolvers.Admin]
	resolvers.HostResolver
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

func SetAppIcon(ctx context.Context, input *AppIconInput) (*SetAppIconOutput, error) {

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
		Location: input.GenerateURL("assets/app_icon.png"),
	}, nil
}
