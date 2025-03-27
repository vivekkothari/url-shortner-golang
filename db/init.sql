CREATE TABLE url_shortner
(
    id               BIGINT      NOT NULL
        PRIMARY KEY,
    long_url         TEXT        NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL,
    last_accessed_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX url_shortner_long_url_uindex ON url_shortner (long_url);
