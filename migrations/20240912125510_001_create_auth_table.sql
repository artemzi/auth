-- +goose Up
CREATE TYPE ROLE AS ENUM ('1', '2');

CREATE TABLE "user" (
    id serial primary key,
    name text not null,
    email text not null,
    password text not null,
    password_confirm text not null,
    role ROLE default '1',
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table "user";
drop type role;
