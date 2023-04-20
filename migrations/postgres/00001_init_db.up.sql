CREATE TABLE links
(
    id SERIAL PRIMARY KEY,
    long TEXT unique,
    short TEXT unique,
    cookie TEXT
);
INSERT INTO links VALUES (0, 'http://localhost:8080/', '0', 'start');