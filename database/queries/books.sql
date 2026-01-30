-- name: InsertBookEvent :exec
INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, published_date, thumbnail_url, occurred_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListBooks :many
SELECT book_id, code, title, authors, publisher, published_date, thumbnail_url, occurred_at
FROM book_events e1
WHERE event_type = 'created'
AND NOT EXISTS (
    SELECT 1 FROM book_events e2
    WHERE e2.book_id = e1.book_id AND e2.event_type = 'deleted'
)
ORDER BY occurred_at DESC
LIMIT ? OFFSET ?;

-- name: CountBooks :one
SELECT COUNT(*) AS cnt
FROM book_events e1
WHERE event_type = 'created'
AND NOT EXISTS (
    SELECT 1 FROM book_events e2
    WHERE e2.book_id = e1.book_id AND e2.event_type = 'deleted'
);

-- name: SearchBooks :many
SELECT book_id, code, title, authors, publisher, published_date, thumbnail_url, occurred_at
FROM book_events e1
WHERE event_type = 'created'
AND NOT EXISTS (
    SELECT 1 FROM book_events e2
    WHERE e2.book_id = e1.book_id AND e2.event_type = 'deleted'
)
AND (title LIKE ? OR authors LIKE ?)
ORDER BY occurred_at DESC
LIMIT ? OFFSET ?;

-- name: CountSearchBooks :one
SELECT COUNT(*) AS cnt
FROM book_events e1
WHERE event_type = 'created'
AND NOT EXISTS (
    SELECT 1 FROM book_events e2
    WHERE e2.book_id = e1.book_id AND e2.event_type = 'deleted'
)
AND (title LIKE ? OR authors LIKE ?);
