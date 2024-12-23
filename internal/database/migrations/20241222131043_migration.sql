-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
create table user_notifications
(
    id         bigint unsigned auto_increment primary key,
    user_id    bigint unsigned not null,
    data       json,
    status     integer         not null default 1,
    updated_at timestamp       not null default now() on update current_timestamp,
    created_at timestamp       not null default now(),
    index idx_u_n_user_id using btree (user_id),
    index idx_u_n_created_at using btree (created_at),
    index idx_u_n_status using hash (status)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE user_notifications;