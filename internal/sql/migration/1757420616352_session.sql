CREATE TABLE session (
    id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    token TEXT NOT NULL,
    csrf_token TEXT NOT NULL,
    expires_at INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    PRIMARY KEY (id),
    UNIQUE (token),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE INDEX idx_session_user_id ON session(user_id);