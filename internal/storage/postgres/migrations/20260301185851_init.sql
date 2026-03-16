-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	login text UNIQUE NOT NULL,
	password text NOT NULL,
	created_at timestamp DEFAULT now()
);

CREATE INDEX CONCURRENTLY ON users (login); 

CREATE TABLE IF NOT EXISTS log_pass_data (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description text,
    login text,
    password text,
	created_at timestamp DEFAULT now(),
    user_id uuid NOT NULL references users(id) ON DELETE CASCADE
);

CREATE INDEX CONCURRENTLY ON log_pass_data (user_id);

CREATE TABLE IF NOT EXISTS text_data (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description text,
    content text,
	created_at timestamp DEFAULT now(),
    user_id uuid NOT NULL references users(id) ON DELETE CASCADE
);

CREATE INDEX CONCURRENTLY ON text_data (user_id);

CREATE TABLE IF NOT EXISTS file_data (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description text,
    name text,
    content bytea,
	created_at timestamp DEFAULT now(),
    user_id uuid NOT NULL references users(id) ON DELETE CASCADE
);

CREATE INDEX CONCURRENTLY ON file_data (user_id);

CREATE TABLE IF NOT EXISTS bank_card_data (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description text,
    holder text,
    number text,
	created_at timestamp DEFAULT now(),
    user_id uuid NOT NULL references users(id) ON DELETE CASCADE
);

CREATE INDEX CONCURRENTLY ON bank_card_data (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS log_pass_data;

DROP TABLE IF EXISTS text_data;

DROP TABLE IF EXISTS file_data;

DROP TABLE IF EXISTS bank_card_data;

DROP TABLE IF EXISTS users;
-- +goose StatementEnd
