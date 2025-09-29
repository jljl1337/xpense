CREATE TABLE expense (
    id TEXT NOT NULL,
    book_id TEXT NOT NULL,
    category_id TEXT NOT NULL,
    payment_method_id TEXT NOT NULL,
    date TEXT NOT NULL,
    amount REAL NOT NULL,
    remark TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE,
    FOREIGN KEY (payment_method_id) REFERENCES payment_method(id) ON DELETE CASCADE,

    CHECK (
        date GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]' AND
        date(date) IS NOT NULL AND
        date(date) = date
    )
);

CREATE INDEX idx_expense_book_id ON expense(book_id);
CREATE INDEX idx_expense_category_id ON expense(category_id);
CREATE INDEX idx_expense_payment_method_id ON expense(payment_method_id);