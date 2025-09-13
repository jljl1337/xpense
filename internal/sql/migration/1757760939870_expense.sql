CREATE TABLE expense (
    id TEXT PRIMARY KEY,
    book_id TEXT NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    category_id TEXT NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    payment_method_id TEXT NOT NULL REFERENCES payment_method(id) ON DELETE CASCADE,
    date INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    remark TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);