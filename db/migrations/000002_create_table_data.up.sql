create table if not exists datas
(
    id          bigserial
    primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    user_id     bigint not null
        constraint fk_datas_users
            references users,
    type        int not null,
    value       varchar not null,
    description varchar
);

create index if not exists idx_datas_user_id
    on datas (user_id);

create index if not exists idx_datas_deleted_at
    on datas (deleted_at);

