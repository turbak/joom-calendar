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
    start_date timestamp NOT NULL,
    is_repeated boolean NOT NULL,
    rrule text,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE INDEX events_start_date_index ON events (start_date);

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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event_attendees;
DROP TABLE events;
DROP TABLE users;
DROP TABLE event_invites;
-- +goose StatementEnd
