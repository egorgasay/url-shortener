CREATE TABLE urls
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    longURL TEXT,
    shortURL TEXT,
    cookie TEXT
);EOQ
INSERT INTO urls (id) VALUES (1);