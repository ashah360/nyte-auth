create extension if not exists "uuid-ossp";

create table if not exists users (
    id uuid default uuid_generate_v4(),
    email text not null,
    password text not null,
    role text not null,
    discord_id text,
    upgraded boolean default False,
    access_expires_at timestamp without time zone default now(),
    created_at timestamp without time zone default now(),

    primary key (id)
);
