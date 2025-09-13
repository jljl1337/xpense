CREATE TABLE payment_method (
    id TEXT PRIMARY KEY,
    book_id TEXT NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE INDEX idx_payment_method_book_id ON payment_method(book_id);