CREATE TABLE urls
(
    id INTEGER primary key autoincrement,
    long TEXT unique,
    short TEXT unique,
    cookie TEXT unique
);
INSERT INTO urls VALUES (0, 'http://localhost:8080/', '0', 'start');