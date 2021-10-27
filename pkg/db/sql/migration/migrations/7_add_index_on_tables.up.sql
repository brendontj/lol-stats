CREATE UNIQUE INDEX idx_schedule_matches_id ON schedule.matches (external_reference);
CREATE INDEX idx_schedule_matches_team_a_name ON schedule.matches (team_a_name);
CREATE INDEX idx_schedule_matches_team_b_name ON schedule.matches (team_b_name);

CREATE UNIQUE INDEX idx_game_games_id ON game.games (id);
CREATE INDEX idx_game_games_match_id ON game.games (matchid);
CREATE INDEX idx_game_games_stats_game_id ON game.games_stats (gameid);
