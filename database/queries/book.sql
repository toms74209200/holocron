-- name: InsertBookEvent :exec
INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, published_date, thumbnail_url, occurred_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetBookByBookId :one
SELECT book_id, code, title, authors, publisher, published_date, thumbnail_url, occurred_at
FROM book_events
WHERE book_id = ? AND event_type = 'created'
LIMIT 1;

-- name: CountBookByBookId :one
SELECT COUNT(*) AS cnt
FROM book_events
WHERE book_id = ? AND event_type = 'created';

-- name: ListBooks :many
SELECT book_id, code, title, authors, publisher, published_date, thumbnail_url, occurred_at
FROM book_events
WHERE event_type = 'created'
ORDER BY occurred_at DESC
LIMIT ? OFFSET ?;

-- name: CountBooks :one
SELECT COUNT(*) AS cnt
FROM book_events
WHERE event_type = 'created';

-- name: GetBookByCode :one
SELECT book_id, code, title, authors, publisher, published_date, thumbnail_url, occurred_at
FROM book_events
WHERE code = ? AND event_type = 'created'
LIMIT 1;
