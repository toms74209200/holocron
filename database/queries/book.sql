-- name: GetBookByBookId :one
SELECT book_id, code, title, authors, publisher, published_date, thumbnail_url, occurred_at
FROM book_events
WHERE book_id = ? AND event_type = 'created'
  AND NOT EXISTS (
    SELECT 1 FROM book_events AS del
    WHERE del.book_id = book_events.book_id AND del.event_type = 'deleted'
  )
LIMIT 1;

-- name: CountBookByBookId :one
SELECT COUNT(*) AS cnt
FROM book_events
WHERE book_id = ? AND event_type = 'created';
