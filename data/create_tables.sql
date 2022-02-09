CREATE TABLE IF NOT EXISTS cookies (
    name  TEXT
        constraint cookies_pk
            primary key,
    value TEXT not null
);
