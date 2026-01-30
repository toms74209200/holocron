-- name: GetBookByBookId :one
SELECT book_id, code, title, authors, publisher, published_date, thumbnail_url, occurred_at
FROM book_events
WHERE book_id = ? AND event_type = 'created'
LIMIT 1;

-- name: CountBookByBookId :one
SELECT COUNT(*) AS cnt
FROM book_events
WHERE book_id = ? AND event_type = 'created';
