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
	ID          uint      `json:"id" db:"id"`
	CreatorID   uint      `json:"creator_id" db:"creator_id"`
	CreatorName string    `json:"creator_name" db:"creator_name"`
	Text        string    `json:"text" db:"text"`
	MediaURL    string    `json:"media_url" db:"media_url"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
}
