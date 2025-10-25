package helper

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	filename := hex.EncodeToString(func() []byte { b := make([]byte, 16); _, _ = rand.Read(b); return b }()) + filepath.Ext(header.Filename)

	destPath := filepath.Join(uploadDir, filename)
	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	return filename, nil
}
