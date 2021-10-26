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

func (s *Storage) GetLastMatchStats(teamName string, numberOfPastGames int, gameTime time.Time, gameMoment int) ([]lol_transformer.MatchGameStats, error) {
	queryGetLastMatchStats := `
SELECT 
       g.id,
       m.external_reference,
       m.team_a_name,
       m.team_b_name,
	   gs.blue_team_barons
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
  AND 
	m.game_time < $2
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
	rows, err := s.pool.Query(context.Background(), queryGetLastMatchStats, teamName, gameTime)
	if err != nil {
		return nil, errors.Wrap(err,"[storage error] unable to get last stats of historic matches")
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

func (s *Storage) UpdateMatchWithPastStats(tx Transaction, gameID uuid.UUID, teamOrder string, stats lol_transformer.StatsInfo, numberOfPastGames, gameMoment int) error {
	queryUpdateStatsTeamALast3At15 :=
		`UPDATE schedule.matches
		 SET 
			team_a_3_gold_total_mean_at15 = $2
			team_a_3_kills_mean_at15 = $3
			team_a_3_inhibitors_mean_at15 = $4
			team_a_3_dragons_mean_at15 = $5
			team_a_3_towers_mean_at15 = $6
		WHERE id = $1;`

	queryUpdateStatsTeamALast3At25 :=
		`UPDATE schedule.matches
		 SET 
			team_a_3_gold_total_mean_at25 = $2
			team_a_3_kills_mean_at25 = $3
			team_a_3_inhibitors_mean_at25 = $4
			team_a_3_dragons_mean_at25 = $5
			team_a_3_towers_mean_at25 = $6
			team_a_3_barons_mean_at25 = $7	
		WHERE id = $1;`

	queryUpdateStatsTeamALast5At15 :=
		`UPDATE schedule.matches
		 SET
			team_a_5_gold_total_mean_at15 = $2
			team_a_5_kills_mean_at15 = $3
			team_a_5_inhibitors_mean_at15 = $4
			team_a_5_dragons_mean_at15 = $5
			team_a_5_towers_mean_at15 = $6
		WHERE id = $1;`

	queryUpdateStatsTeamALast5At25 :=
		`UPDATE schedule.matches
		 SET 
			team_a_5_gold_total_mean_at25 = $2
			team_a_5_kills_mean_at25 = $3
			team_a_5_inhibitors_mean_at25 = $4
			team_a_5_dragons_mean_at25 = $5
			team_a_5_towers_mean_at25 = $6
			team_a_5_barons_mean_at25 = $7	
		WHERE id = $1;`

	queryUpdateStatsTeamBLast3At15 :=
		`UPDATE schedule.matches
		 SET 
			team_b_3_gold_total_mean_at15 = $2
			team_b_3_kills_mean_at15 = $3
			team_b_3_inhibitors_mean_at15 = $4
			team_b_3_dragons_mean_at15 = $5
			team_b_3_towers_mean_at15 = $6
		WHERE id = $1;`

	queryUpdateStatsTeamBLast3At25 :=
		`UPDATE schedule.matches
		 SET 
			team_b_3_gold_total_mean_at25 = $2
			team_b_3_kills_mean_at25 = $3
			team_b_3_inhibitors_mean_at25 = $4
			team_b_3_dragons_mean_at25 = $5
			team_b_3_towers_mean_at25 = $6
			team_b_3_barons_mean_at25 = $7	
		WHERE id = $1;`

	queryUpdateStatsTeamBLast5At15 :=
		`UPDATE schedule.matches
		 SET 
			team_b_5_gold_total_mean_at15 = $2
			team_b_5_kills_mean_at15 = $3
			team_b_5_inhibitors_mean_at15 = $4
			team_b_5_dragons_mean_at15 = $5
			team_b_5_towers_mean_at15 = $6
		WHERE id = $1;`

	queryUpdateStatsTeamBLast5At25 :=
		`UPDATE schedule.matches
		 SET 
			team_b_5_gold_total_mean_at25 = $2
			team_b_5_kills_mean_at25 = $3
			team_b_5_inhibitors_mean_at25 = $4
			team_b_5_dragons_mean_at25 = $5
			team_b_5_towers_mean_at25 = $6
			team_b_5_barons_mean_at25 = $7	
		WHERE id = $1;`

	var query string
	if teamOrder == "A" {
		if numberOfPastGames == 3 {
			if gameMoment == 4 {
				query = queryUpdateStatsTeamALast3At15
			}

			if gameMoment == 6 {
				query = queryUpdateStatsTeamALast3At25
			}
		}

		if numberOfPastGames == 5 {
			if gameMoment == 4 {
				query = queryUpdateStatsTeamALast5At15
			}

			if gameMoment == 6 {
				query = queryUpdateStatsTeamALast5At25
			}
		}
	}

	if teamOrder == "B" {
		if numberOfPastGames == 3 {
			if gameMoment == 4 {
				query = queryUpdateStatsTeamBLast3At15
			}

			if gameMoment == 6 {
				query = queryUpdateStatsTeamBLast3At25
			}
		}

		if numberOfPastGames == 5 {
			if gameMoment == 4 {
				query = queryUpdateStatsTeamBLast5At15
			}

			if gameMoment == 6 {
				query = queryUpdateStatsTeamBLast5At25
			}
		}
	}

	if gameMoment == 6 {
		_, err := tx.Exec(
			context.Background(),
			query,
			gameID,
			stats.NumberOfTotalGoldMean,
			stats.NumberOfKillsMean,
			stats.NumberOfInhibitorsMean,
			stats.NumberOfDragonsMean,
			stats.NumberOfTowersMean,
			stats.NumberOfBaronsMean)
		if err != nil {
			return errors.Wrapf(err, fmt.Sprintf("[storage error] unable to store winner ratio for game with external reference = %s", gameID))
		}
	}

	if gameMoment == 4 {
		_, err := tx.Exec(
			context.Background(),
			query,
			gameID,
			stats.NumberOfTotalGoldMean,
			stats.NumberOfKillsMean,
			stats.NumberOfInhibitorsMean,
			stats.NumberOfDragonsMean,
			stats.NumberOfTowersMean)
		if err != nil {
			return errors.Wrapf(err, fmt.Sprintf("[storage error] unable to store winner ratio for game with external reference = %s", gameID))
		}
	}

	return nil
}