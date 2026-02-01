-- name: InsertLendingEvent :exec
INSERT INTO lending_events (
    event_id,
    lending_id,
    book_id,
    borrower_id,
    event_type,
    due_date,
    occurred_at
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetCurrentLending :one
SELECT
    lending_id,
    book_id,
    borrower_id,
    due_date,
    occurred_at as borrowed_at
FROM lending_events
WHERE book_id = ?
    AND event_type = 'borrowed'
    AND NOT EXISTS (
        SELECT 1
        FROM lending_events returned
        WHERE returned.lending_id = lending_events.lending_id
            AND returned.event_type = 'returned'
    )
ORDER BY occurred_at DESC
LIMIT 1;

-- name: IsBookBorrowed :one
SELECT EXISTS(
    SELECT 1
    FROM lending_events
    WHERE book_id = ?
        AND event_type = 'borrowed'
        AND NOT EXISTS (
            SELECT 1
            FROM lending_events returned
            WHERE returned.lending_id = lending_events.lending_id
                AND returned.event_type = 'returned'
        )
) as is_borrowed;

-- name: GetLatestDueDate :one
SELECT due_date
FROM lending_events
WHERE lending_id = ?
    AND event_type IN ('borrowed', 'due_date_extended')
    AND due_date IS NOT NULL
ORDER BY occurred_at DESC
LIMIT 1;
