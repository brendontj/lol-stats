package services

import (
	"context"
	"fmt"
	lol_transformer "github.com/brendontj/lol-stats/pkg/lol-transformer"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
}

type Storage struct {
	pool *pgxpool.Pool
}

func (s *Storage) GetAllMatchesWithoutProcessedData() ([]lol_transformer.Match, error) {
	queryGetAllMatchesWithoutProcessedData := `
SELECT 
	id,
	team_a_name,
	team_b_name,
	event_start_time
FROM schedule.matches as m
WHERE 
	m.team_a_5_form_ratio = NULL
AND 
	m.team_b_5_form_ratio = NULL
ORDER BY
	m.event_start_time ASC;
`
	rows, err := s.pool.Query(context.Background(), queryGetAllMatchesWithoutProcessedData)
	if err != nil {
		return nil, errors.Wrap(err,"[storage error] unable to get all matches")
	}
	defer rows.Close()

	var matches []lol_transformer.Match
	for rows.Next() {
		var m lol_transformer.Match
		err = rows.Scan(
			&m.ID,
			&m.TeamAName,
			&m.TeamBName,
			&m.GameTime)
		if err != nil {
			return nil, errors.Wrapf(err, "[storage error] unable to scan match with external reference = %s", m.ID)
		}
		matches = append(matches, m)
	}

	return matches, nil
}

func (s *Storage) GetLastMatchResults(tx Transaction, teamName string, numberOfPastGames int, gameTime time.Time) ([]lol_transformer.MatchResult, error) {
	queryGetLastMatchResults := `
SELECT
	team_a_name,
	team_b_name,
	team_a_game_wins,
	team_b_game_wins,
	bestOf
FROM schedule.matches as m
WHERE 
	(m.team_a_name = $1 OR m.team_b_name = $1)
AND
	m.event_start_time < $2
ORDER BY
	m.event_start_time DESC
LIMIT $3;`

	rows, err := tx.Query(
		context.Background(),
		queryGetLastMatchResults,
		teamName,
		gameTime,
		numberOfPastGames)
	if err != nil {
		return nil, errors.Wrap(err,"[storage error] unable to get last match results")
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

func (s *Storage) UpdateMatchWithWinnerRatio(tx Transaction, gameID uuid.UUID, teamOrder string, ratio float64, numberOfGames int) error {
	queryUpdateMatchWithWinnerRatioTeamALast3 :=
		`UPDATE schedule.matches
		 SET team_a_3_form_ratio = $2   
		WHERE id = $1;`

	queryUpdateMatchWithWinnerRatioTeamALast5 :=
		`UPDATE schedule.matches
		 SET team_a_5_form_ratio = $2   
		WHERE id = $1;`

	queryUpdateMatchWithWinnerRatioTeamBLast3 :=
		`UPDATE schedule.matches
		 SET team_b_3_form_ratio = $2   
		WHERE id = $1;`

	queryUpdateMatchWithWinnerRatioTeamBLast5 :=
		`UPDATE schedule.matches
		 SET team_b_5_form_ratio = $2   
		WHERE id = $1;`

	var query string
	if teamOrder == "A" {
		if numberOfGames == 3 {
			query = queryUpdateMatchWithWinnerRatioTeamALast3
		}

		if numberOfGames == 5 {
			query = queryUpdateMatchWithWinnerRatioTeamALast5
		}
	}

	if teamOrder == "B" {
		if numberOfGames == 3 {
			query = queryUpdateMatchWithWinnerRatioTeamBLast3
		}

		if numberOfGames == 5 {
			query = queryUpdateMatchWithWinnerRatioTeamBLast5
		}
	}
	_, err := tx.Exec(
		context.Background(),
		query,
		gameID,
		ratio)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("[storage error] unable to store winner ratio for game with external reference = %s", gameID))
	}

	return nil
}

func (s *Storage) GetLastMatchResults(tx Transaction, teamName string, numberOfPastGames int, gameTime time.Time) ([]lol_transformer.MatchGameStats, error) {
	queryGetLastMatchStats := `
SELECT 
	id,
	team_a_name,
	team_b_name,
	event_start_time
FROM schedule.matches as m
WHERE 
	m.team_a_5_form_ratio = NULL
AND 
	m.team_b_5_form_ratio = NULL
ORDER BY
	m.event_start_time ASC;
`
	rows, err := s.pool.Query(context.Background(), queryGetAllMatchesWithoutProcessedData)
	if err != nil {
		return nil, errors.Wrap(err,"[storage error] unable to get all matches")
	}
	defer rows.Close()

	var matches []lol_transformer.Match
	for rows.Next() {
		var m lol_transformer.Match
		err = rows.Scan(
			&m.ID,
			&m.TeamAName,
			&m.TeamBName,
			&m.GameTime)
		if err != nil {
			return nil, errors.Wrapf(err, "[storage error] unable to scan match with external reference = %s", m.ID)
		}
		matches = append(matches, m)
	}

	return matches, nil
}