CREATE TABLE users (
    id          CHAR(36) PRIMARY KEY,
    username    VARCHAR(255) NOT NULL UNIQUE,
    email       VARCHAR(255) UNIQUE,
    full_name   TEXT NOT NULL,
    avatar_key  TEXT,
    create_time TIMESTAMP(6) NOT NULL,
    update_time TIMESTAMP(6) NOT NULL
) ENGINE = InnoDB;

CREATE TABLE local_auth_strategies (
    user_id       CHAR(36) PRIMARY KEY,
    password_hash TEXT NOT NULL,

    CONSTRAINT fk_local_auth_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE google_auth_strategies (
    user_id CHAR(36) PRIMARY KEY,
    sub     TEXT NOT NULL,

    CONSTRAINT fk_google_auth_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE outbox_messages (
    id            CHAR(36) PRIMARY KEY,
    partition_key TEXT,
    source        TEXT NOT NULL,
    type          TEXT NOT NULL,
    payload       JSON NOT NULL,
    occur_time    TIMESTAMP(6) NOT NULL
) ENGINE = InnoDB;

CREATE TABLE sessions (
    id          CHAR(36) PRIMARY KEY,
    user_id     CHAR(36) NOT NULL,
    token       BINARY(32) NOT NULL UNIQUE,
    expire_time TIMESTAMP(6) NOT NULL,
    revoke_time TIMESTAMP(6),

    INDEX idx_sessions_user_id (user_id),

    CONSTRAINT fk_session_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE = InnoDB;