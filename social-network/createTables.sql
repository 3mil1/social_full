CREATE TABLE IF NOT EXISTS user
(
    id            VARCHAR(36)  NOT NULL PRIMARY KEY,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name    VARCHAR(255) NOT NULL,
    last_name     VARCHAR(255) NOT NULL,
    birthday      date         not null,
    image         VARCHAR(255),
    nickname      VARCHAR(255),
    about         TEXT,
    created_at    timestamp    not null default (datetime('now','localtime')),
    updated_at    timestamp    default 0,
    is_private    BOOLEAN NOT NULL CHECK (is_private IN (0, 1))
);

CREATE TABLE IF NOT EXISTS post
(
    id         INTEGER                 NOT NULL,
    user_id    CHAR(36)                NOT NULL,
    title      VARCHAR(255)            NOT NULL default '',
    content    TEXT                    NOT NULL,
    image      TEXT,
    parent_id  INTEGER NULL,
    created_at timestamp not null default (datetime('now','localtime')),
    updated_at timestamp default 0,
    privacy     INT     NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user (id),
    FOREIGN KEY (parent_id) REFERENCES post (id),
    PRIMARY KEY (id AUTOINCREMENT)

);

CREATE TABLE IF NOT EXISTS post_access
(
    id      INTEGER  NOT NULL,
    post_id INTEGER  NOT NULL,
    user_id CHAR(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user (id),
    FOREIGN KEY (post_id) REFERENCES post (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS follower
(
    id        INTEGER     NOT NULL,
    source_id VARCHAR(36) NOT NULL,
    target_id VARCHAR(36) NOT NULL,
    status    Integer     NOT NULL,
    FOREIGN KEY (source_id) REFERENCES user (id),
    FOREIGN KEY (target_id) REFERENCES user (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS chat
(
    id       INTEGER     NOT NULL,
    user1_id VARCHAR(36) NOT NULL,
    user2_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (user1_id) REFERENCES user (id),
    FOREIGN KEY (user2_id) REFERENCES user (id),
    PRIMARY KEY (id AUTOINCREMENT),
    UNIQUE (user1_id, user2_id)
);

CREATE TABLE IF NOT EXISTS chat_message
(
    id         INTEGER     NOT NULL,
    chat_id    INTEGER     NOT NULL,
    user_id    VARCHAR(36) NOT NULL,
    content    TEXT        NOT NULL,
    created_at  timestamp    not null default (datetime('now','localtime')),
    FOREIGN KEY (user_id) REFERENCES user (id),
    FOREIGN KEY (chat_id) REFERENCES chat (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS message_status
(
    message_id INTEGER     NOT NULL,
    user_id    VARCHAR(36) NOT NULL,
    read       BOOLEAN     NOT NULL DEFAULT false,
    FOREIGN KEY (message_id) REFERENCES chat_message (id),
    FOREIGN KEY (user_id) REFERENCES user (id)
);
CREATE TABLE IF NOT EXISTS sessions
(
    refresh_token VARCHAR(255) NOT NULL UNIQUE,
    user_id       VARCHAR(36)  NOT NULL,
    device        varchar(255) not null,
    ip            varchar TEXT not null,
    PRIMARY KEY (user_id, device, ip) ON CONFLICT REPLACE,
    FOREIGN KEY (user_id) REFERENCES user (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS groups
(
    id          INTEGER      NOT NULL,
    created_by  VARCHAR(36)  NOT NULL,
    title       VARCHAR(255) NOT NULL UNIQUE,
    description TEXT         NOT NULL,
    FOREIGN KEY (created_by) REFERENCES user (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS group_post
(
    id         INTEGER      NOT NULL,
    group_id   INTEGER      NOT NULL,
    user_id    VARCHAR(36)  NOT NULL,
    title      VARCHAR(255) NOT NULL,
    content    TEXT         NOT NULL,
    image      VARCHAR(255),
    parent_id  INTEGER,
    created_at timestamp    not null default (datetime('now', 'localtime')),
    updated_at INTEGER,
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (user_id) REFERENCES user (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS group_message
(
    id         INTEGER     NOT NULL,
    group_id   INTEGER     NOT NULL,
    user_id    VARCHAR(36) NOT NULL,
    content    TEXT        NOT NULL,
    created_at timestamp   not null default (datetime('now', 'localtime')),
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (user_id) REFERENCES user (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS group_message_status
(
    message_id INTEGER     NOT NULL,
    user_id    VARCHAR(36) NOT NULL,
    read       BOOLEAN     NOT NULL,
    FOREIGN KEY (message_id) REFERENCES group_message (id),
    FOREIGN KEY (user_id) REFERENCES user (id)
);

CREATE TABLE IF NOT EXISTS group_member
(
    id       INTEGER     NOT NULL,
    group_id INTEGER     NOT NULL,
    user_id  VARCHAR(36) NOT NULL,
    role     VARCHAR(36) NOT NULL,
    status   INTEGER     NOT NULL,

    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (user_id) REFERENCES user (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS group_event
(
    id          INTEGER      NOT NULL,
    group_id    INTEGER      NOT NULL,
    user_id     VARCHAR(36)  NOT NULL,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    event_date  timestamp    not null default (datetime('now','localtime')),
    created_at  timestamp    not null default (datetime('now','localtime')),

    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (user_id) REFERENCES user (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS event_participants
(
    id       INTEGER     NOT NULL,
    event_id INTEGER     NOT NULL,
    user_id  VARCHAR(36) NOT NULL,
    option   INTEGER     NOT NULL,
    FOREIGN KEY (event_id) REFERENCES group_event (id),
    FOREIGN KEY (user_id) REFERENCES user (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS notification_type
(
    id   INTEGER      NOT NULL,
    type VARCHAR(255) NOT NULL,
    PRIMARY KEY (id AUTOINCREMENT)
);
insert or ignore into notification_type (id, type)
VALUES (1, 'group invitation'),
       (2, 'new group member request'),
       (3, 'new event'),
       (4, 'friend request'),
       (5, 'new private message'),
       (6, 'new comment to post'),
       (7, 'group access opened'),
       (8, 'new message in group chat');

CREATE TABLE IF NOT EXISTS notification_obj
(
    id                INTEGER     NOT NULL,
    notification_type INTEGER     NOT NULL,
    object_id         Integer     NOT NULL,
    actor_id          VARCHAR(36) NOT NULL,
    created_at        timestamp   not null default (datetime('now', 'localtime')),
    FOREIGN KEY (notification_type) REFERENCES notification_type (id),
    FOREIGN KEY (actor_id) REFERENCES user (id),
    PRIMARY KEY (id AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS notification
(
    id              INTEGER     NOT NULL,
    receiver_id     VARCHAR(36) NOT NULL,
    notification_id INTEGER     NOT NULL,
    seen            INTEGER     NOT NULL,

    FOREIGN KEY (receiver_id) REFERENCES user (id),
    FOREIGN KEY (notification_id) REFERENCES notification_obj (id),
    PRIMARY KEY (id AUTOINCREMENT)
);