CREATE TABLE payment_method (
    id TEXT NOT NULL,
    book_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE
);

CREATE INDEX idx_payment_method_book_id ON payment_method(book_id);