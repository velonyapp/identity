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

CREATE TABLE outbox_events (
    id             CHAR(36) PRIMARY KEY,
    aggregate_id   TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    type           TEXT NOT NULL,
    payload        JSON NOT NULL,
    occur_time     TIMESTAMP(6) NOT NULL
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

CREATE TABLE email_change_requests (
    id          CHAR(36) PRIMARY KEY,
    token_hash  BINARY(32) NOT NULL UNIQUE,
    user_id     CHAR(36) NOT NULL,
    new_email   VARCHAR(255) NOT NULL,
    create_time TIMESTAMP(6) NOT NULL,
    expire_time TIMESTAMP(6) NOT NULL,
    confirm_time TIMESTAMP(6),

    INDEX idx_change_email_requests_user_id (user_id),

    CONSTRAINT fk_change_email_request_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE = InnoDB;
