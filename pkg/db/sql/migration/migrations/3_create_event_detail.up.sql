CREATE TABLE schedule.events_detail (
    id UUID PRIMARY KEY,
    game_ref VARCHAR NOT NULL,
    event_external_ref  VARCHAR NOT NULL REFERENCES schedule.matches(external_reference),
    tournament_external_ref  VARCHAR NOT NULL,
    league_external_ref VARCHAR NOT NULL
);

CREATE TABLE schedule.events_games (
    event_external_ref  VARCHAR NOT NULL REFERENCES schedule.matches(external_reference),
    game_external_ref VARCHAR NOT NULL,
    game_number INTEGER NOT NULL,
    status  VARCHAR NOT NULL,
    team_a_external_ref  VARCHAR NOT NULL,
    team_b_external_ref VARCHAR NOT NULL,
    team_a_side VARCHAR NOT NULL,
    team_b_side VARCHAR NOT NULL,
    UNIQUE (event_external_ref, game_external_ref)
    );
