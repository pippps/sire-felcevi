-- +goose Up
CREATE TABLE images_md(
    id UUID PRIMARY KEY NOT NULL, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT,
    image_url TEXT NOT NULL,
    description TEXT,
    user_id UUID NOT NULL,
    constraint fk_user_id
	FOREIGN KEY(user_id) 
	REFERENCES users(id)
	ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE images_md;