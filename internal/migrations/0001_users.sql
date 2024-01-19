-- +goose Up
create table if not exists public.users
(
    id serial not null primary key,
    name varchar(100) not null,
    surname varchar(100) not null,
    patronymic varchar(100) not null,
    age integer,
    gender varchar(1),
    nation varchar(2)
);

-- +goose Down
drop table public.users;