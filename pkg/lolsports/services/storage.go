package services

import (
	"context"
	"fmt"
	lol_transformer "github.com/brendontj/lol-stats/pkg/lol-transformer"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	pool *pgxpool.Pool
}

func (s *Storage) ExistsLeague(id string) (bool, error) {
	query := `SELECT EXISTS(SELECT * FROM league.leagues WHERE external_reference=$1)`
	row := s.pool.QueryRow(context.Background(), query, id)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) ExistsEvent(id string) (bool, error) {
	query := `SELECT EXISTS(SELECT * FROM schedule.matches WHERE external_reference=$1)`
	row := s.pool.QueryRow(context.Background(), query, id)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) ExistsEventExternalRef(gameReference string) (bool, error) {
	query := `SELECT EXISTS(SELECT * FROM schedule.events_detail WHERE game_ref=$1)`
	row := s.pool.QueryRow(context.Background(), query, gameReference)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) ExistsGameID(gameID string) (bool, error) {
	query := `SELECT EXISTS(SELECT * FROM game.games WHERE gameID=$1)`
	row := s.pool.QueryRow(context.Background(), query, gameID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) InsertLeagues(leagueData *lolsports.LeagueData) error {
	queryInsertLeagueMetadata :=
		`INSERT INTO league.leagues (ID, external_reference, slug, name, region, image, priority)
	VALUES ($1, $2, $3, $4, $5, $6, $7);`

	for _, league := range leagueData.ScheduleContent.Leagues {
		leagueWasPersisted, err := s.ExistsLeague(league.ID)
		if err != nil {
			return err
		}

		if leagueWasPersisted {
			continue
		}

		_, err = s.pool.Exec(
			context.Background(),
			queryInsertLeagueMetadata,
			uuid.NewV4(),
			league.ID,
			league.Slug,
			league.Name,
			league.Region,
			league.Image,
			league.Priority)
		if err != nil {
			return errors.Wrapf(err, fmt.Sprintf("[storage error] unable to store league with id = %s", league.ID))
		}
		fmt.Printf("Inserted league (%s, %s, %s, %s, %s, %d) into database\n",
			league.ID,
			league.Slug,
			league.Name,
			league.Region,
			league.Image,
			league.Priority)
	}
	return nil
}

func (s *Storage) SaveEvent(event lolsports.Events) error {
	queryInsertMatchMetadata :=
		`INSERT INTO schedule.matches 
(ID, external_reference, team_a_name, team_a_code, team_a_image, team_b_name, team_b_code, team_b_image, 
team_a_record_wins, team_a_record_losses, team_b_record_wins, team_b_record_losses, team_a_game_wins, team_b_game_wins, best_of, event_start_time,
state, league_name) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);`

	matchID := uuid.NewV4()
	_, err := s.pool.Exec(
		context.Background(),
		queryInsertMatchMetadata,
		matchID,
		event.Match.ID,
		event.Match.Teams[0].Name,
		event.Match.Teams[0].Code,
		event.Match.Teams[0].Image,
		event.Match.Teams[1].Name,
		event.Match.Teams[1].Code,
		event.Match.Teams[1].Image,
		event.Match.Teams[0].Record.Wins,
		event.Match.Teams[0].Record.Losses,
		event.Match.Teams[1].Record.Wins,
		event.Match.Teams[1].Record.Losses,
		event.Match.Teams[0].Result.GameWins,
		event.Match.Teams[1].Result.GameWins,
		event.Match.Strategy.Count,
		event.StartTime,
		event.State,
		event.League.Name)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("[storage error] unable to store match with id = %s", event.Match.ID))
	}
	return nil
}

func (s *Storage) SaveEventDetail(eventDetail *lolsports.EventDetailData, gameRef string) error {
	queryInsertEventDetailMetadata :=
		`INSERT INTO schedule.events_detail (ID, game_ref, event_external_ref, tournament_external_ref, league_external_ref) 
VALUES ($1, $2, $3, $4, $5);`

	_, err := s.pool.Exec(
		context.Background(),
		queryInsertEventDetailMetadata,
		uuid.NewV4(),
		gameRef,
		eventDetail.Data.Event.ID,
		eventDetail.Data.Event.Tournament.TournamentID,
		eventDetail.Data.Event.League.ID)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("[storage error] unable to store event with external reference = %s", eventDetail.Data.Event.ID))
	}
	return nil
}

func (s *Storage) SaveGameInfo(eventDetail *lolsports.EventDetailData) error {
	queryInsertGameInfoMetadata :=
		`INSERT INTO schedule.events_games (event_external_ref, game_external_ref, game_number, status, team_a_external_ref, team_b_external_ref, team_a_side, team_b_side) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	for _, game := range eventDetail.Data.Event.Match.Games {
		_, err := s.pool.Exec(
			context.Background(),
			queryInsertGameInfoMetadata,
			eventDetail.Data.Event.ID,
			game.ID,
			game.Number,
			game.State,
			game.Teams[0].TeamID,
			game.Teams[1].TeamID,
			game.Teams[0].Side,
			game.Teams[1].Side)
		if err != nil {
			return errors.Wrapf(err, fmt.Sprintf("[storage error] unable to store game info with external reference = %s", game.ID))
		}
	}
	return nil
}

func (s *Storage) GetLeagues() ([]lolsports.League, error) {
	queryGetAllLeagues := `
SELECT 
	external_reference,
	slug,
	name,
	region,
	image,
	priority
FROM league.leagues;
`
	rows, err := s.pool.Query(context.Background(), queryGetAllLeagues)
	if err != nil {
		return nil, errors.Wrap(err, "[storage error] unable to get league from storage")
	}
	defer rows.Close()

	var leagues []lolsports.League
	for rows.Next() {
		var l lolsports.League
		err = rows.Scan(
			&l.ID,
			&l.Slug,
			&l.Name,
			&l.Region,
			&l.Image,
			&l.Priority)
		if err != nil {
			return nil, errors.Wrapf(err, "[storage error] unable to scan league with external reference = %s", l.ID)
		}
		leagues = append(leagues, l)
	}

	return leagues, nil
}

func (s *Storage) GetLastMatchResults(teamName string, numberOfPastGames int) ([]lol_transformer.MatchResult, error) {
	queryGetLastMatchResults := `
SELECT
	team_a_name,
	team_b_name,
	team_a_game_wins,
	team_b_game_wins,
	best_of
FROM schedule.matches as m
WHERE (m.team_a_name = $1 OR m.team_b_name = $1)
ORDER BY m.event_start_time DESC
LIMIT $2;`

	rows, err := s.pool.Query(
		context.Background(),
		queryGetLastMatchResults,
		teamName,
		numberOfPastGames)
	if err != nil {
		return nil, errors.Wrap(err, "[storage error] unable to get last match results")
	}
	defer rows.Close()

	var matchesResults []lol_transformer.MatchResult
	for rows.Next() {
		var m lol_transformer.MatchResult
		err = rows.Scan(
			&m.TeamAName,
			&m.TeamBName,
			&m.TeamAGameWins,
			&m.TeamBGameWins,
			&m.BestOf,
		)
		if err != nil {
			return nil, err
		}
		matchesResults = append(matchesResults, m)
	}

	return matchesResults, nil
}

func (s *Storage) GetLastMatchStats(teamName string, numberOfPastGames int, gameMoment int) ([]lol_transformer.MatchGameStats, error) {
	queryGetLastMatchStats := `
SELECT 
       g.id,
       m.external_reference,
       m.team_a_name,
       m.team_b_name,
	   gs.blue_team_barons,
       gs.blue_team_dragons,
       gs.blue_team_inhibitors,
       gs.blue_team_total_gold,
       gs.blue_team_total_kills,
       gs.blue_team_towers,
       gs.red_team_barons,
       gs.red_team_dragons,
       gs.red_team_inhibitors,
       gs.red_team_total_gold,
       gs.red_team_total_kills,
       gs.red_team_towers
FROM schedule.matches as m
         INNER JOIN game.games g
                    ON m.external_reference = g.matchid
         INNER JOIN game.games_stats gs
                    ON g.id = gs.gameid
WHERE
    (m.team_a_name = $1 OR m.team_b_name = $1)
  AND
    m.team_a_5_gold_total_mean_at15 is NULL
  AND
    m.team_a_3_gold_total_mean_at15 IS NULL
  AND
    m.team_b_5_gold_total_mean_at15 IS NULL
  AND
    m.team_b_3_gold_total_mean_at15 IS NULL
GROUP BY g.id,
         m.event_start_time,
         gs.timestamp,m.external_reference,
         m.team_a_name,
         m.team_b_name,
         gs.red_team_barons,
         gs.red_team_dragons,
         gs.red_team_inhibitors,
         gs.red_team_total_gold,
         gs.red_team_total_kills,
         gs.red_team_towers,
         gs.blue_team_towers,
         gs.blue_team_total_kills,
         gs.blue_team_total_gold,
         gs.blue_team_inhibitors,
         gs.blue_team_dragons,
         gs.blue_team_barons
ORDER BY m.event_start_time DESC, gs.timestamp ASC
`
	rows, err := s.pool.Query(context.Background(), queryGetLastMatchStats, teamName)
	if err != nil {
		return nil, errors.Wrap(err, "[storage error] unable to get last stats of historic matches")
	}
	defer rows.Close()

	var matches []lol_transformer.MatchGameStats
	var currentGameID uuid.UUID
	count := 0
	for rows.Next() {
		var m lol_transformer.MatchGameStats
		err = rows.Scan(
			&m.GameID,
			&m.MatchExternalReference,
			&m.TeamAName,
			&m.TeamBName,
			&m.TeamBlueTotalBarons,
			&m.TeamBlueDragons,
			&m.TeamBlueTotalInhibitors,
			&m.TeamBlueTotalGold,
			&m.TeamBlueTotalKills,
			&m.TeamBlueTotalTowers,
			&m.TeamRedTotalBarons,
			&m.TeamRedDragons,
			&m.TeamRedTotalInhibitors,
			&m.TeamRedTotalGold,
			&m.TeamRedTotalKills,
			&m.TeamRedTotalTowers)
		if err != nil {
			return nil, errors.Wrapf(err, "[storage error] unable to scan match with external reference = %s", m.MatchExternalReference)
		}

		if currentGameID == uuid.Nil {
			currentGameID = m.GameID
		}

		if m.GameID != currentGameID {
			currentGameID = m.GameID
			count = 0
		}

		if count == gameMoment {
			matches = append(matches, m)
		}

		if len(matches) == numberOfPastGames {
			break
		}

		count += 1
	}

	return matches, nil
}
