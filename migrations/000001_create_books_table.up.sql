create table if not exists books (
    id bigserial primary key,
    title varchar not null,
    author_first_name varchar not null,
    author_last_name varchar not null,
    genre varchar not null,
    quantity int not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);