CREATE TABLE users (
    id          CHAR(36) PRIMARY KEY,
    username    VARCHAR(255) NOT NULL,
    email       VARCHAR(255),
    avatar_id   CHAR(36),
    create_time TIMESTAMP(6) NOT NULL,
    version     INT NOT NULL,

    UNIQUE KEY uq_users_username (username),
    UNIQUE KEY uq_users_email    (email)
) ENGINE = InnoDB;

CREATE TABLE avatars (
    id         CHAR(36) PRIMARY KEY,
    user_id    CHAR(36) NOT NULL,
    object_key VARCHAR(128) NOT NULL,

    CONSTRAINT fk_avatars_user_id
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE = InnoDB;

ALTER TABLE users
    ADD CONSTRAINT fk_users_avatar_id
        FOREIGN KEY (avatar_id)
        REFERENCES avatars (id)
        ON DELETE SET NULL;

CREATE TABLE email_change_requests (
    user_id     CHAR(36) PRIMARY KEY,
    token       TEXT NOT NULL,
    email       VARCHAR(255) NOT NULL,
    create_time TIMESTAMP(6) NOT NULL,
    expire_time TIMESTAMP(6) NOT NULL,

    CONSTRAINT fk_email_change_requests_user_id
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE,

    INDEX idx_email_change_requests_expire_time (expire_time)
) ENGINE = InnoDB;

CREATE TABLE avatar_change_requests (
    user_id     CHAR(36) PRIMARY KEY,
    avatar_id   CHAR(36) NOT NULL,
    create_time TIMESTAMP(6) NOT NULL,
    expire_time TIMESTAMP(6) NOT NULL,

    CONSTRAINT fk_avatar_change_requests_user_id
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE,

    INDEX idx_avatar_change_requests_expire_time (expire_time)
) ENGINE = InnoDB;

CREATE TABLE local_auth_strategies (
    user_id       CHAR(36) PRIMARY KEY,
    password_hash TEXT NOT NULL,

    CONSTRAINT fk_local_auth_strategies_user_id
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE sessions (
    id          CHAR(36) PRIMARY KEY,
    user_id     CHAR(36) NOT NULL,
    token_hash  BINARY(32) NOT NULL,
    create_time TIMESTAMP(6) NOT NULL,
    expire_time TIMESTAMP(6) NOT NULL,
    revoke_time TIMESTAMP(6),
    version     INT NOT NULL,

    UNIQUE KEY uq_sessions_token_hash (token_hash),

    CONSTRAINT fk_sessions_user_id
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE outbox_events (
    id             CHAR(36) PRIMARY KEY,
    type           TEXT NOT NULL,
    aggregate_id   TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    payload        JSON NOT NULL,
    occur_time     DATETIME(6) NOT NULL,

    INDEX idx_outbox_events_occur_time (occur_time)
) ENGINE = InnoDB;

CREATE TABLE outbox_jobs (
    id           CHAR(36) PRIMARY KEY,
    type         TEXT NOT NULL,
    payload      JSON NOT NULL,
    enqueue_time DATETIME(6) NOT NULL,

    INDEX idx_outbox_jobs_enqueue_time (enqueue_time)
) ENGINE = InnoDB;
