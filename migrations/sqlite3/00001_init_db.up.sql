CREATE TABLE links
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    long TEXT,
    short TEXT,
    cookie TEXT
);
INSERT INTO links VALUES (0, 'http://localhost:8080/', '0', 'start');