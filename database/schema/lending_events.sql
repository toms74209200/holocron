CREATE TABLE lending_events (
    event_id TEXT PRIMARY KEY,
    lending_id TEXT NOT NULL,
    book_id TEXT NOT NULL,
    borrower_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    due_date TEXT,
    occurred_at TEXT NOT NULL
);

CREATE INDEX idx_lending_events_lending_id ON lending_events(lending_id);
CREATE INDEX idx_lending_events_book_id ON lending_events(book_id);
CREATE INDEX idx_lending_events_borrower_id ON lending_events(borrower_id);
