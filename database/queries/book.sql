-- name: GetBookByBookId :one
SELECT
    b.book_id,
    b.code,
    b.title,
    b.authors,
    b.publisher,
    b.published_date,
    b.thumbnail_url,
    b.occurred_at
FROM book_events b
INNER JOIN (
    SELECT book_id, MAX(occurred_at) as latest_occurred_at
    FROM book_events
    WHERE event_type IN ('created', 'deleted')
    GROUP BY book_id
) latest ON b.book_id = latest.book_id
        AND b.occurred_at = latest.latest_occurred_at
WHERE b.book_id = ?
  AND b.event_type = 'created';

-- name: CountBookByBookId :one
SELECT COUNT(*) AS cnt
FROM book_events
WHERE book_id = ? AND event_type = 'created';

-- name: InsertBookUpdateEvent :exec
INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, published_date, thumbnail_url, occurred_at)
VALUES (?, ?, 'updated', ?, ?, ?, ?, ?, ?, ?);
