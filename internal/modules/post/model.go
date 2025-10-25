package post

import "time"

type Post struct {
	ID          uint      `db:"id"`
	CreatorID   uint      `db:"creator_id"`
	Text        string    `db:"text"`
	MediaURL    string    `db:"media_url"`
	PublishedAt time.Time `db:"published_at"`
}
