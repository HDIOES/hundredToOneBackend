
-- +migrate Up
CREATE TABLE GAMES (
    ID SERIAL PRIMARY KEY,
    BODY BYTEA
);
-- +migrate Down
