-- +goose Up
CREATE TYPE ROLE AS ENUM ('ROLE_USER', 'ROLE_ADMIN');

CREATE TABLE "user" (
    id serial primary key,
    name text not null,
    email text not null,
    password text not null,
    password_confirm text not null,
    role ROLE default 'ROLE_USER',
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table "user";
drop type role;
