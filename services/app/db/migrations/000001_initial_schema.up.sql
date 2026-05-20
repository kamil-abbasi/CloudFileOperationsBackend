CREATE TABLE IF NOT EXISTS users (
	id UUID NOT NULL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	max_storage BIGINT NOT NULL,
  storage_used BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS directories (
	id UUID NOT NULL PRIMARY KEY,
	user_id UUID NOT NULL,
	parent_id UUID,
	name TEXT NOT NULL,
	location TEXT NOT NULL,
	FOREIGN KEY (parent_id)
		REFERENCES directories(id)
		ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS files (
	id UUID NOT NULL PRIMARY KEY,
	user_id UUID NOT NULL,
	directory_id UUID,
	name TEXT NOT NULL,
	size INTEGER NOT NULL,
	location TEXT NOT NULL,
	FOREIGN KEY(directory_id)
		REFERENCES directories(id)
		ON DELETE CASCADE
);

CREATE TYPE item_type AS ENUM('file', 'directory');

CREATE TABLE IF NOT EXISTS directory_items (
	id UUID NOT NULL PRIMARY KEY,
	user_id UUID NOT NULL,
	type item_type NOT NULL
);