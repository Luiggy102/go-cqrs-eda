DROP TABLE if exists feeds;

CREATE TABLE feeds (
    id varchar(32) PRIMARY KEY,
    title varchar(255) not null,
    description varchar(255) not null,
    created_at timestamp not null default now()
);
