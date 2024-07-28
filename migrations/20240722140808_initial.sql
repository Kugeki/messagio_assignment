-- +goose Up
-- +goose StatementBegin
create table messages
(
    id        serial
        constraint messages_pk
            primary key,
    content   varchar            not null,
    processed bool default false not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages;
-- +goose StatementEnd
