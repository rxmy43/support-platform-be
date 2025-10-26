package helper

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/rxmy43/support-platform/internal/config"
)

func SaveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	cfg := config.Load()

	cld, err := cloudinary.NewFromParams(
		cfg.Cloudinary.Name,
		cfg.Cloudinary.ApiKey,
		cfg.Cloudinary.ApiSecret,
	)
	if err != nil {
		return "", err
	}

	resp, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		ResourceType: "auto",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
