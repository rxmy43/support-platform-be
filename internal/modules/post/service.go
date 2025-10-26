package post

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/config"
	"github.com/rxmy43/support-platform/internal/helper"
	"github.com/rxmy43/support-platform/internal/modules/user"
)

type PostService struct {
	postRepo *PostRepo
	userRepo *user.UserRepo
}

func NewPostService(postRepo *PostRepo, userRepo *user.UserRepo) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

const maxFileSize = 5 * 1024 * 1024

var allowedExts = map[string]bool{
	".jpg":  true,
	".png":  true,
	".jpeg": true,
}

func (s *PostService) Create(ctx context.Context, req PostCreateRequest) *apperror.AppError {
	var fieldErrs []apperror.FieldError

	if req.Text == "" {
		fieldErrs = append(fieldErrs, apperror.NewFieldError("text", apperror.CodeFieldRequired))
	}

	if _, err := s.userRepo.FindByID(ctx, req.CreatorID); err != nil {
		if err == sql.ErrNoRows {
			return apperror.Unauthorized("creator id not found", apperror.CodeUnauthorizedOperation)
		}
		return apperror.InternalServer("failed executing find user by creator id").WithCause(err)
	}

	if req.Header.Size > maxFileSize {
		fieldErrs = append(fieldErrs, apperror.NewFieldError("file", apperror.CodeFileTooLarge))
	}

	ext := filepath.Ext(req.Header.Filename)
	if !allowedExts[ext] {
		fieldErrs = append(fieldErrs, apperror.NewFieldError("file", apperror.CodeFileTypeInvalid))
	}

	if len(fieldErrs) > 0 {
		return apperror.ValidationError("create post validation error", fieldErrs)
	}

	secureURL, err := helper.SaveUploadedFile(req.File, req.Header)
	if err != nil {
		return apperror.InternalServer("failed when uploading file").WithCause(err)
	}

	newPost := &Post{
		CreatorID: req.CreatorID,
		Text:      req.Text,
		MediaURL:  secureURL,
	}

	if err := s.postRepo.Create(ctx, newPost); err != nil {
		return apperror.InternalServer("failed when creating new post").WithCause(err)
	}

	return nil
}

func (s *PostService) GenerateCaption(ctx context.Context, tone string) (string, *apperror.AppError) {
	cfg := config.Load()
	apiKey := cfg.GroqAPIKey

	url := "https://api.groq.com/openai/v1/chat/completions"

	if tone == "" {
		tone = "friendly and engaging"
	}

	prompt := fmt.Sprintf(
		`You are an AI assistant helping creators write short, emotionally engaging captions 
for their posts on a creator support platform (like Patreon, Trakteer, or Saweria). 
Use a %s tone. The caption should sound natural, personal, and subtly encourage appreciation or donations.
Output ONLY the caption text — no explanations, no quotes, no introductions, no phrases like "Here’s your caption."
The caption must not exceed 500 characters.`,
		tone,
	)

	payload := map[string]interface{}{
		"model": "llama-3.3-70b-versatile",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a creative writing assistant for a content creator support platform.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.85,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", apperror.InternalServer("failed to connect to Groq API").WithCause(err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", apperror.InternalServer("failed to send request to Groq").WithCause(err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", apperror.InternalServer("failed to read Groq response").WithCause(err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", apperror.Wrap(apperror.CodeExternalAPIRequestFailed, resp.StatusCode, "Groq API error", errors.New(string(resBody)))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(resBody, &result); err != nil {
		return "", apperror.InternalServer("failed to parse Groq response").WithCause(err)
	}

	if len(result.Choices) == 0 {
		return "", apperror.InternalServer("failed generated caption").WithCause(err)
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

func (s *PostService) FindAll(ctx context.Context, cursor, userID *uint) ([]PostResponse, *uint, *apperror.AppError) {
	posts, nextCursor, err := s.postRepo.GetPosts(ctx, cursor, userID)
	if err != nil {
		return []PostResponse{}, nil, apperror.InternalServer("failed get all posts").WithCause(err)
	}
	return posts, nextCursor, nil
}
