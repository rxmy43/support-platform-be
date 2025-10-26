package post

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/repo"
)

type PostRepo struct {
	*repo.BaseRepo[Post]
}

func NewPostRepo(DB *sqlx.DB) *PostRepo {
	return &PostRepo{
		BaseRepo: &repo.BaseRepo[Post]{
			DB:        DB,
			TableName: "posts",
		},
	}
}

func (r *PostRepo) GetPosts(ctx context.Context, cursor, userID *uint) ([]PostResponse, *uint, error) {
	posts := []PostResponse{}
	var err error

	queryBase := `
		SELECT p.id, p.creator_id, u.name AS creator_name, p.text, p.media_url, p.published_at
		FROM posts p
		JOIN users u ON u.id = p.creator_id
	`

	args := []any{}
	where := ""

	// filter user
	if userID != nil {
		where = "WHERE u.id = $1"
		args = append(args, *userID)
	}

	// filter cursor
	if cursor != nil {
		if where == "" {
			where = "WHERE p.id < $1"
			args = append(args, *cursor)
		} else {
			where += " AND p.id < $2"
			args = append(args, *cursor)
		}
	}

	query := fmt.Sprintf(`
		%s
		%s
		ORDER BY p.id DESC
		LIMIT 10
	`, queryBase, where)

	err = r.DB.SelectContext(ctx, &posts, query, args...)
	if err != nil {
		return nil, nil, err
	}

	var nextCursor *uint
	if len(posts) > 0 {
		nextCursor = &posts[len(posts)-1].ID
	}

	return posts, nextCursor, nil
}
