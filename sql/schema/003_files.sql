-- +goose Up
CREATE TABLE files(
    id UUID PRIMARY KEY UNIQUE NOT NULL, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size BIGINT NOT NULL,
    etag TEXT NOT NULL,
    user_id UUID NOT NULL,
    constraint fk_user_id
	FOREIGN KEY(user_id) 
	REFERENCES users(id)
	ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE files;