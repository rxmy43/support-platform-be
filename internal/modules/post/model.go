package post

type Post struct {
	ID        uint   `db:"id"`
	CreatorID uint   `db:"creator_id"`
	Text      string `db:"text"`
	MediaURL  string `db:"media_url"`
}
