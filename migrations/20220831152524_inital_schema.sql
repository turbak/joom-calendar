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
    event_id integer NOT NULL,
    user_id integer NOT NULL,
    status TEXT NOT NULL default 'unconfirmed',
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    
    PRIMARY KEY (event_id, user_id)
);

CREATE TABLE event_invites (
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
    day_of_week TEXT NOT NULL,
    day_of_month TEXT NOT NULL,
    month_of_year TEXT NOT NULL,
    week_of_month TEXT NOT NULL,
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
