CREATE TABLE memes (
    id SERIAL PRIMARY KEY,
    phash TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL,
    source TEXT NOT NULL,
    source_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);