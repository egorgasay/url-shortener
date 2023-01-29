CREATE TABLE urls
(
    id SERIAL PRIMARY KEY,
    long TEXT unique,
    short TEXT unique,
    cookie TEXT
);EOQ
INSERT INTO urls VALUES (0, 'http://localhost:8080/', '0', 'start');