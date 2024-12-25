CREATE TABLE users (
    id bigserial not null primary key,
    email varchar not null unique,
    login varchar not null,
    password varchar not null,
    date TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Добавляем таблицу posts
CREATE TABLE posts (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content VARCHAR(300) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE admins (
    id bigserial not null primary key,
    email varchar not null unique,
    name varchar not null,
    password varchar not null,
    date TIMESTAMP NOT NULL DEFAULT NOW()
);
