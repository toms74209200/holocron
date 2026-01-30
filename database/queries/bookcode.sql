-- name: InsertBookEvent :exec
INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, published_date, thumbnail_url, occurred_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetBookByCode :one
SELECT book_id, code, title, authors, publisher, published_date, thumbnail_url, occurred_at
FROM book_events
WHERE code = ? AND event_type = 'created'
LIMIT 1;
