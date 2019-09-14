CREATE TABLE IF NOT EXISTS users (
    id integer NOT NULL,
    fac text NOT NULL,
    groupnum text NOT NULL,
    notifications boolean NOT NULL,
    ignorelist text[]
);