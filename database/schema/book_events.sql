CREATE TABLE book_events (
    event_id TEXT PRIMARY KEY,
    book_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    code TEXT,
    title TEXT,
    authors TEXT,
    publisher TEXT,
    published_date TEXT,
    thumbnail_url TEXT,
    delete_reason TEXT,
    delete_memo TEXT,
    occurred_at TEXT NOT NULL
);

CREATE INDEX idx_book_events_book_id ON book_events(book_id);
