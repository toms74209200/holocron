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

-- name: ListBorrowingBooksByBorrowerID :many
WITH my_current_lendings AS (
    SELECT le.book_id, le.lending_id, le.occurred_at as borrowed_at
    FROM lending_events le
    WHERE le.borrower_id = ?
        AND le.event_type = 'borrowed'
        AND NOT EXISTS (
            SELECT 1
            FROM lending_events returned
            WHERE returned.lending_id = le.lending_id
                AND returned.event_type = 'returned'
        )
),
deleted_books AS (
    SELECT book_id, MAX(occurred_at) as deleted_at
    FROM book_events
    WHERE event_type = 'deleted'
    GROUP BY book_id
),
latest_books AS (
    SELECT
        e1.book_id,
        e1.code,
        e1.title,
        e1.authors,
        e1.publisher,
        e1.published_date,
        e1.thumbnail_url,
        (SELECT MIN(e_created.occurred_at)
         FROM book_events e_created
         WHERE e_created.book_id = e1.book_id
           AND e_created.event_type = 'created'
           AND e_created.occurred_at > COALESCE(
               (SELECT MAX(e_del.occurred_at) FROM book_events e_del WHERE e_del.book_id = e1.book_id AND e_del.event_type = 'deleted'),
               '1970-01-01T00:00:00Z'
           )
        ) as created_at,
        ROW_NUMBER() OVER (PARTITION BY e1.book_id ORDER BY e1.occurred_at DESC) as rn
    FROM book_events e1
    LEFT JOIN deleted_books d ON e1.book_id = d.book_id
    WHERE e1.event_type IN ('created', 'updated')
        AND (d.deleted_at IS NULL OR e1.occurred_at > d.deleted_at)
)
SELECT
    lb.book_id,
    lb.code,
    lb.title,
    lb.authors,
    lb.publisher,
    lb.published_date,
    lb.thumbnail_url,
    lb.created_at,
    mcl.borrowed_at,
    (
        SELECT ldd.due_date
        FROM lending_events ldd
        WHERE ldd.lending_id = mcl.lending_id
            AND ldd.event_type IN ('borrowed', 'due_date_extended')
            AND ldd.due_date IS NOT NULL
        ORDER BY ldd.occurred_at DESC
        LIMIT 1
    ) as due_date
FROM my_current_lendings mcl
JOIN latest_books lb ON lb.book_id = mcl.book_id AND lb.rn = 1
ORDER BY mcl.borrowed_at DESC;
