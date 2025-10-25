package post

import (
	"mime/multipart"
	"time"
)

type PostCreateRequest struct {
	CreatorID uint
	Text      string
	File      multipart.File
	Header    *multipart.FileHeader
}

type PostResponse struct {
	ID          uint      `json:"id"`
	CreatorID   uint      `json:"creator_id"`
	CreatorName string    `json:"creator_name"`
	Text        string    `json:"text"`
	MediaURL    string    `json:"media_url"`
	PublishedAt time.Time `json:"published_at"`
}
