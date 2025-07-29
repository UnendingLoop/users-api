CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY SERIAL,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS friends(
    requester INTEGER NOT NULL,
    accepter INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (requester, accepter),
    FOREIGN KEY (requester) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (accepter) REFERENCES users(id) ON DELETE CASCADE
);