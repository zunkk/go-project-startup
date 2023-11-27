create table if not exists "user"
(
    id         int       not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp,

    constraint user_pk primary key (id)
);
