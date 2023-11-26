create table if not exists "user"
(
    id          int       not null,
    create_time timestamp not null,
    update_time timestamp not null,
    delete_time timestamp,

    constraint user_pk primary key (id)
);
