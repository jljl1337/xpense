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

CREATE INDEX idx_expense_book_id ON expense(book_id);
CREATE INDEX idx_expense_category_id ON expense(category_id);
CREATE INDEX idx_expense_payment_method_id ON expense(payment_method_id);