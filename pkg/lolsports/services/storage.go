package services

import (
	"context"
	"fmt"
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
	row := s.pool.QueryRow(context.Background(),query,id)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) ExistsEvent(id string) (bool, error) {
	query := `SELECT EXISTS(SELECT * FROM schedule.matches WHERE external_reference=$1)`
	row := s.pool.QueryRow(context.Background(),query,id)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) InsertLeagues(leagueData *lolsports.LeagueData) error{
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

func (s *Storage) SaveEventDetail(eventDetail *lolsports.EventDetailData) error {
	queryInsertEventDetailMetadata :=
		`INSERT INTO schedule.events_detail (ID, event_external_ref, tournament_external_ref, league_external_ref) 
VALUES ($1, $2, $3, $4);`

	_, err := s.pool.Exec(
		context.Background(),
		queryInsertEventDetailMetadata,
		uuid.NewV4(),
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

func (s *Storage) GetLeagues() ([]lolsports.League,error) {
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
		return nil, errors.Wrap(err,"[storage error] unable to get league from storage")
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