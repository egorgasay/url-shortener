CREATE TABLE urls
(
    id INTEGER primary key autoincrement,
    long TEXT unique,
    short TEXT unique
);
INSERT INTO urls VALUES (0, 'http://localhost:8080/', '0');