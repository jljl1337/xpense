CREATE TABLE user (
    id TEXT NOT NULL,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    PRIMARY KEY (id),
    UNIQUE (username)
);

CREATE INDEX idx_user_username ON user(username);