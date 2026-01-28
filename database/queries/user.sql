-- name: InsertUserEvent :exec
INSERT INTO user_events (event_id, user_id, event_type, name, occurred_at)
VALUES (?, ?, ?, ?, ?);

-- name: GetUserByUserId :one
SELECT user_id, name, occurred_at
FROM user_events
WHERE user_id = ? AND event_type = 'created'
LIMIT 1;

-- name: CountUserByUserId :one
SELECT COUNT(*) AS cnt
FROM user_events
WHERE user_id = ? AND event_type = 'created';
