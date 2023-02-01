CREATE TABLE urls
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    long TEXT,
    short TEXT,
    cookie TEXT
);
INSERT INTO urls VALUES (0, 'http://localhost:8080/', '0', 'start');