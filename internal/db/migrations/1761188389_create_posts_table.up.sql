CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    creator_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    media_url VARCHAR(500),
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_posts_creator_id ON posts(creator_id);
