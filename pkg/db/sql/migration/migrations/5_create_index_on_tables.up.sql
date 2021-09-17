CREATE INDEX participant_stats_game_id
    ON game.participants_stats (gameid);

CREATE INDEX participant_info_game_id
    ON game.participants_info (gameid);

CREATE INDEX game_stats_game_id
    ON game.games_stats (gameid);

CREATE INDEX games_id
    ON game.games (id);

CREATE INDEX matches_external_reference
    ON schedule.matches (external_reference);

CREATE INDEX events_games_external_ref
    ON schedule.events_games (game_external_ref);

CREATE INDEX events_details_external_ref
    ON schedule.events_detail (event_external_ref);