CREATE TABLE user_events (
    event_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    name TEXT NOT NULL,
    occurred_at TEXT NOT NULL
);

CREATE INDEX idx_user_events_user_id ON user_events(user_id);
