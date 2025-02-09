-- +goose Up
-- +goose StatementBegin
create table users
(
	id           bigint unsigned auto_increment primary key,
	first_name   varchar(100)      not null,
	second_name  varchar(100)      not null,
	email        varchar(100)      not null,
	phone_number varchar(100)      not null,
	password     varchar(255)      not null,
	status       smallint unsigned not null,
	hash         varchar(128)      null,
	roles        json              null,
	confirmed_at timestamp         null,
	updated_at   timestamp         not null,
	created_at   timestamp         not null,
	constraint users_email_uniq unique (email),
	constraint uses_phone_number_uniq unique (phone_number)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
