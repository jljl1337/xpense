CREATE TABLE category (
    id TEXT NOT NULL,
    book_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE
);

CREATE INDEX idx_category_book_id ON category(book_id);