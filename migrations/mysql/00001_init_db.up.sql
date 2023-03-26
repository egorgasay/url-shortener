CREATE TABLE urls
(
    id INT AUTO_INCREMENT,
    longURL TEXT,
    shortURL TEXT,
    cookie TEXT,
    PRIMARY KEY (id)
);
INSERT INTO urls (longURL, shortURL, cookie) VALUES ('http://localhost:8080/', '0', 'start');