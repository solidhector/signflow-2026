-- name: CreateDocument :one
INSERT INTO documents (owner_id, filename, content_type, size, file_hash, storage_key)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: ListDocumentsByOwner :many
SELECT * FROM documents
WHERE owner_id = $1
ORDER BY created_at DESC;