CREATE SCHEMA game;

CREATE TABLE game.games (
    id UUID NOT NULL PRIMARY KEY,
    gameID varchar not null,
    matchID varchar not null,
    patch_version varchar not null,
    blueTeamID varchar not null,
    redTeamID varchar not null
);

CREATE TABLE game.participants_stats (
    gameID uuid not null references game.games(id),
    game_externalID varchar not null,
    participantID int not null,
    game_timestamp timestamp not null,
    level int not null,
    kills int not null,
    deaths int not null,
    assists int not null,
    total_gold_earned int not null,
    creep_score int not null,
    kill_participation decimal(15,2) not null,
    champion_damage_share decimal(15,2) not null,
    wards_placed int not null,
    wards_destroyed int not null,
    unique (game_timestamp,participantID,game_externalID)
);

CREATE TABLE game.participants_info (
  gameID uuid not null references game.games(id),
  participantID int not null,
  championID varchar not null,
  esportsPlayerID varchar not null,
  summonerName varchar not null,
  role varchar not null
);

CREATE TABLE game.games_stats (
  gameID uuid not null references game.games(id),
  timestamp timestamp not null,
  gameState varchar not null,
  blueTeamID varchar not null,
  redTeamID varchar not null,
  blue_team_total_gold int not null,
  blue_team_inhibitors int not null,
  blue_team_towers int not null,
  blue_team_barons int not null,
  blue_team_total_kills int not null,
  blue_team_dragons TEXT [] not null,
  red_team_total_gold int not null,
  red_team_inhibitors int not null,
  red_team_towers int not null,
  red_team_barons int not null,
  red_team_total_kills int not null,
  red_team_dragons TEXT [] not null
);
