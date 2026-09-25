CREATE TABLE users (
    id          CHAR(36) PRIMARY KEY,
    username    VARCHAR(255) NOT NULL UNIQUE,
    full_name   TEXT NOT NULL,
    email       VARCHAR(255) UNIQUE,
    avatar_key  VARCHAR(128),
    create_time TIMESTAMP(6) NOT NULL,
    update_time TIMESTAMP(6) NOT NULL
) ENGINE = InnoDB;

CREATE TABLE email_change_requests (
    user_id     CHAR(36) PRIMARY KEY,
    value       VARCHAR(255) NOT NULL,
    time        TIMESTAMP(6) NOT NULL,
    expire_time TIMESTAMP(6) NOT NULL,

    CONSTRAINT fk_change_email_request_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE avatar_change_requests (
    user_id     CHAR(36) PRIMARY KEY,
    value       VARCHAR(128) NOT NULL,
    time        TIMESTAMP(6) NOT NULL,
    expire_time TIMESTAMP(6) NOT NULL,

    CONSTRAINT fk_change_avatar_request_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
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
