-- 1. User table
CREATE TABLE users (
    id          UUID            NOT NULL,
    email       VARCHAR(255)    NOT NULL,
    name        VARCHAR(255),
    picture     TEXT,
    created_at  TIMESTAMP       DEFAULT NOW(),
    updated_at  TIMESTAMP       DEFAULT NOW(),

    PRIMARY KEY (id),
    UNIQUE (email)
);

-- 2. User Identities table
CREATE TABLE user_identities (
    id                  UUID            NOT NULL,
    user_id             UUID            NOT NULL,
    provider            VARCHAR(50)     NOT NULL, -- e.g. 'google', 'github'
    provider_user_id    VARCHAR(255)    NOT NULL, -- e.g. Google sub
    access_token        TEXT,
    refresh_token       TEXT,
    expires_at          TIMESTAMP,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT fk_user_identities FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (provider, provider_user_id)
);
