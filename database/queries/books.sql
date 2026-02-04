-- name: InsertBookEvent :exec
INSERT INTO book_events (event_id, book_id, event_type, code, title, authors, publisher, published_date, thumbnail_url, occurred_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListBooks :many
WITH deleted_books AS (
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
        e1.occurred_at as updated_at,
        ROW_NUMBER() OVER (PARTITION BY e1.book_id ORDER BY e1.occurred_at DESC) as rn
    FROM book_events e1
    LEFT JOIN deleted_books d ON e1.book_id = d.book_id
    WHERE e1.event_type IN ('created', 'updated')
        AND (d.deleted_at IS NULL OR e1.occurred_at > d.deleted_at)
),
current_lendings AS (
    SELECT
        le.book_id,
        le.borrower_id,
        ue.name as borrower_name,
        le.occurred_at as borrowed_at,
        ROW_NUMBER() OVER (PARTITION BY le.book_id ORDER BY le.occurred_at DESC) as rn
    FROM lending_events le
    INNER JOIN user_events ue ON ue.user_id = le.borrower_id AND ue.event_type = 'created'
    WHERE le.event_type = 'borrowed'
        AND le.occurred_at > COALESCE(
            (SELECT MAX(e2.occurred_at) FROM book_events e2 WHERE e2.book_id = le.book_id AND e2.event_type = 'deleted'),
            '1970-01-01T00:00:00Z'
        )
        AND NOT EXISTS (
            SELECT 1
            FROM lending_events returned
            WHERE returned.lending_id = le.lending_id
                AND returned.event_type = 'returned'
        )
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
    lb.updated_at,
    cl.borrower_id,
    cl.borrower_name,
    cl.borrowed_at
FROM latest_books lb
LEFT JOIN current_lendings cl ON lb.book_id = cl.book_id AND cl.rn = 1
WHERE lb.rn = 1
ORDER BY lb.updated_at DESC
LIMIT ? OFFSET ?;

-- name: CountBooks :one
SELECT COUNT(DISTINCT e1.book_id) AS cnt
FROM book_events e1
WHERE e1.event_type IN ('created', 'updated')
    AND e1.occurred_at > COALESCE(
        (SELECT MAX(e2.occurred_at) FROM book_events e2 WHERE e2.book_id = e1.book_id AND e2.event_type = 'deleted'),
        '1970-01-01T00:00:00Z'
    );

-- name: SearchBooks :many
WITH deleted_books AS (
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
        e1.occurred_at as updated_at,
        ROW_NUMBER() OVER (PARTITION BY e1.book_id ORDER BY e1.occurred_at DESC) as rn
    FROM book_events e1
    LEFT JOIN deleted_books d ON e1.book_id = d.book_id
    WHERE e1.event_type IN ('created', 'updated')
        AND (d.deleted_at IS NULL OR e1.occurred_at > d.deleted_at)
        AND (e1.title LIKE ? OR e1.authors LIKE ?)
),
current_lendings AS (
    SELECT
        le.book_id,
        le.borrower_id,
        ue.name as borrower_name,
        le.occurred_at as borrowed_at,
        ROW_NUMBER() OVER (PARTITION BY le.book_id ORDER BY le.occurred_at DESC) as rn
    FROM lending_events le
    INNER JOIN user_events ue ON ue.user_id = le.borrower_id AND ue.event_type = 'created'
    WHERE le.event_type = 'borrowed'
        AND le.occurred_at > COALESCE(
            (SELECT MAX(e2.occurred_at) FROM book_events e2 WHERE e2.book_id = le.book_id AND e2.event_type = 'deleted'),
            '1970-01-01T00:00:00Z'
        )
        AND NOT EXISTS (
            SELECT 1
            FROM lending_events returned
            WHERE returned.lending_id = le.lending_id
                AND returned.event_type = 'returned'
        )
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
    lb.updated_at,
    cl.borrower_id,
    cl.borrower_name,
    cl.borrowed_at
FROM latest_books lb
LEFT JOIN current_lendings cl ON lb.book_id = cl.book_id AND cl.rn = 1
WHERE lb.rn = 1
ORDER BY lb.updated_at DESC
LIMIT ? OFFSET ?;

-- name: CountSearchBooks :one
SELECT COUNT(DISTINCT e1.book_id) AS cnt
FROM book_events e1
WHERE e1.event_type IN ('created', 'updated')
    AND e1.occurred_at > COALESCE(
        (SELECT MAX(e2.occurred_at) FROM book_events e2 WHERE e2.book_id = e1.book_id AND e2.event_type = 'deleted'),
        '1970-01-01T00:00:00Z'
    )
    AND (e1.title LIKE ? OR e1.authors LIKE ?);
