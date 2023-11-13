CREATE TABLE IF NOT EXISTS urls
(
    url TEXT PRIMARY KEY,
    description TEXT,
    created_at TEXT,
    read_at TEXT,
    priority INT
);