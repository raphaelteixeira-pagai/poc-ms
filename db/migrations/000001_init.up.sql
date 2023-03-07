CREATE TABLE IF NOT EXISTS wallet (
    id SERIAL PRIMARY KEY,
    balance integer,
    owner varchar(36)
);