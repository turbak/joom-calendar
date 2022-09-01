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
    title text NOT NULL,
    description text NOT NULL,
    duration integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE event_attendees (
    id serial PRIMARY KEY,
    event_id integer NOT NULL,
    user_id integer NOT NULL,
    status TEXT NOT NULL default 'pending',
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE event_repeats (
    id serial PRIMARY KEY,
    event_id integer NOT NULL,
    repeat_start_date timestamp NOT NULL,
    days_of_week integer ARRAY,
    day_of_month integer,
    month_of_year integer,
    week_of_month TEXT,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event_repeats;
DROP TABLE event_attendees;
DROP TABLE events;
DROP TABLE users;
-- +goose StatementEnd
