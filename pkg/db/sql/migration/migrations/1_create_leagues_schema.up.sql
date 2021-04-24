CREATE SCHEMA league;

CREATE TABLE league.leagues (
    ID UUID NOT NULL PRIMARY KEY,
    external_reference VARCHAR NOT NULL,
    slug VARCHAR NOT NULL,
    name VARCHAR NOT NULL ,
    region VARCHAR NOT NULL,
    image VARCHAR NOT NULL,
    priority INTEGER NOT NULL
);
