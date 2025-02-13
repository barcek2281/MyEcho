CREATE TABLE users (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    email varchar(255) not null unique,
    login varchar not null,
    password varchar not null,
    date TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN DEFAULT false,
    is_prime BOOLEAN DEFAULT false
);

CREATE TABLE barcode (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    barcode INT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- Добавляем таблицу posts
CREATE TABLE posts (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    like_count INT,
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

CREATE TABLE messages (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    sender VARCHAR NOT NULL,
    receiver VARCHAR NOT NULL,
    msg VARCHAR NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE allow (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    emailFirst VARCHAR NOT NULL,
    emailSecond VARCHAR NOT NULL
);
