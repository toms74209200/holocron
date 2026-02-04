-- name: GetBookByBookId :one
SELECT
    e1.book_id,
    e1.code,
    e1.title,
    e1.authors,
    e1.publisher,
    e1.published_date,
    e1.thumbnail_url,
    e1.occurred_at
FROM book_events e1
WHERE e1.book_id = ?
    AND e1.event_type IN ('created', 'updated')
    AND e1.occurred_at > COALESCE(
        (SELECT MAX(e2.occurred_at) FROM book_events e2 WHERE e2.book_id = e1.book_id AND e2.event_type = 'deleted'),
        '1970-01-01T00:00:00Z'
    )
ORDER BY e1.occurred_at DESC
LIMIT 1;

-- name: CountBookByBookId :one
SELECT COUNT(*) AS cnt
FROM book_events
WHERE book_id = ? AND event_type = 'created';

-- name: InsertBookUpdateEvent :exec
INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, published_date, thumbnail_url, occurred_at)
VALUES (?, ?, 'updated', ?, ?, ?, ?, ?, ?, ?);

-- name: GetBookBorrowerInfo :one
SELECT
    le.lending_id,
    le.borrower_id,
    ue.name as borrower_name,
    le.occurred_at as borrowed_at
FROM lending_events le
INNER JOIN user_events ue ON ue.user_id = le.borrower_id AND ue.event_type = 'created'
WHERE le.book_id = ?
    AND le.event_type = 'borrowed'
    AND NOT EXISTS (
        SELECT 1
        FROM lending_events returned
        WHERE returned.lending_id = le.lending_id
            AND returned.event_type = 'returned'
    )
ORDER BY le.occurred_at DESC
LIMIT 1;
