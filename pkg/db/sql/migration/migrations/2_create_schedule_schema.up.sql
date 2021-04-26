CREATE SCHEMA schedule;

CREATE TABLE schedule.matches (
    ID UUID NOT NULL PRIMARY KEY,
    external_reference VARCHAR NOT NULL,
    team_a_name VARCHAR NOT NULL,
    team_a_code VARCHAR NOT NULL,
    team_a_image VARCHAR NOT NULL,
    team_b_name VARCHAR NOT NULL,
    team_b_code VARCHAR NOT NULL,
    team_b_image VARCHAR NOT NULL,
    team_a_record_wins INTEGER NOT NULL,
    team_a_record_losses INTEGER NOT NULL,
    team_b_record_wins INTEGER NOT NULL,
    team_b_record_losses INTEGER NOT NULL,
    team_a_game_wins INTEGER NOT NULL,
    team_b_game_wins INTEGER NOT NULL,
    best_of INTEGER NOT NULL
);
