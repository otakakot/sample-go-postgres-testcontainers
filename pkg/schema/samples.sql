-- name: SelectSample :many
SELECT * FROM samples;

-- name: FindSample :one
SELECT * FROM samples WHERE id = $1;

-- name: InsertSample :one
INSERT INTO samples (name) VALUES ($1) RETURNING *;
