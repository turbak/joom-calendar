-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id serial PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email_uindex ON users (email);

CREATE TABLE events (
    id serial PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE event_attendees (
    id serial PRIMARY KEY,
    event_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE event_invites (
    id serial PRIMARY KEY,
    event_id integer NOT NULL,
    user_id integer NOT NULL,
    status TEXT NOT NULL default 'pending',
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE event_repeat (
    id serial PRIMARY KEY,
    event_id integer NOT NULL,
    repeat_start_date timestamp NOT NULL,
    repeat_end_date timestamp,
    days_of_week TEXT ARRAY,
    day_of_month integer,
    month_of_year integer,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event_repeat;
DROP TABLE event_invites;
DROP TABLE event_attendees;
DROP TABLE events;
DROP TABLE users;
-- +goose StatementEnd
