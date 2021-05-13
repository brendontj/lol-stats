package lolsports

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	PopulateLeagues() error
	GetLeagues() ([]League,error)
	GetEventsExternalRef() ([]*string, error)
	PopulateDBScheduleOfLeague(leagueExternalReference string) error
}

type lolService struct {
	storage          *pgxpool.Pool
	esportsApiClient EsportsAPIScrapper
	feedApiClient    FeedAPIScrapper
}

func NewLolService(pgStorage *pgxpool.Pool, esportsApiClient EsportsAPIScrapper, feedApiAclient FeedAPIScrapper) Service {
	return &lolService{
		storage:          pgStorage,
		esportsApiClient: esportsApiClient,
		feedApiClient:    feedApiAclient,
	}
}

func (l *lolService) PopulateLeagues() error {
	leagueData, err := l.esportsApiClient.GetLeagues("pt-BR")
	if err != nil {
		return errors.Wrap(err, "[service error] unable to get leagues")
	}

	queryInsertLeagueMetadata :=
`INSERT INTO league.leagues (ID, external_reference, slug, name, region, image, priority) 
VALUES ($1, $2, $3, $4, $5, $6, $7);`

	for _, league := range leagueData.ScheduleContent.Leagues {
		_, err := l.storage.Exec(
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
			return errors.Wrapf(err, fmt.Sprintf("unable to store league with id = %s", league.ID))
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

func (l *lolService) GetLeagues() ([]League,error) {
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
	rows, err := l.storage.Query(context.Background(), queryGetAllLeagues)
	if err != nil {
		return nil, errors.Wrap(err,"unable to get league from storage")
	}
	defer rows.Close()

	var leagues []League
	for rows.Next() {
		var l League
		err = rows.Scan(
			&l.ID,
			&l.Slug,
			&l.Name,
			&l.Region,
			&l.Image,
			&l.Priority)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to scan league with external reference = %s", l.ID)
		}
		leagues = append(leagues, l)
	}

	return leagues, nil
}

func (l *lolService) PopulateDBScheduleOfLeague(leagueExternalReference string) error {
	olderPage, err := l.populateWithMostRecentScheduleByLeague(leagueExternalReference)
	if err != nil {
		return errors.Wrap(err, "unable to populate with the most recent schedule")
	}

	if olderPage != EmptyField {
		for {
			scheduleData, err := l.esportsApiClient.GetSchedule("pt-BR", leagueExternalReference, olderPage)
			if err != nil {
				return errors.Wrapf(err, "[service error] unable to get schedules with older page = %s",olderPage)
			}
			op, err := l.saveScheduleContent(scheduleData.Data.Schedule)
			if err != nil {
				return err
			}

			if op == EmptyField {
				break
			}
			olderPage = op
		}
	}

	return nil
}

func (l *lolService) populateWithMostRecentScheduleByLeague(leagueExternalReference string) (string, error) {
	scheduleData, err := l.esportsApiClient.GetSchedule("pt-BR", leagueExternalReference, EmptyField)
	if err != nil {
		return EmptyField, errors.Wrap(err, "[service error] unable to get leagues")
	}

	s, err := l.saveScheduleContent(scheduleData.Data.Schedule)
	if err != nil {
		return s, err
	}
	return scheduleData.Data.Schedule.Pages.Older, nil
}

func (l *lolService) saveScheduleContent(scheduleContent Schedule) (string, error) {
	for _, event := range scheduleContent.Events {
		err := l.saveEvent(event)
		if err != nil {
			return EmptyField, err
		}
	}
	return scheduleContent.Pages.Older, nil
}

func (l *lolService) saveEvent(event Events) error {
	queryInsertMatchMetadata :=
		`INSERT INTO schedule.matches 
(ID, external_reference, team_a_name, team_a_code, team_a_image, team_b_name, team_b_code, team_b_image, 
team_a_record_wins, team_a_record_losses, team_b_record_wins, team_b_record_losses, team_a_game_wins, team_b_game_wins, best_of, event_start_time,
state, league_name) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);`

	matchID := uuid.NewV4()
	_, err := l.storage.Exec(
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
		return errors.Wrapf(err, fmt.Sprintf("unable to store match with id = %s", event.Match.ID))
	}
	return nil
}

func (l *lolService) GetEventsExternalRef() ([]*string,error) {
	queryGetAllEventsRef := `
SELECT 
	e.external_reference
FROM schedule.matches e
where e.state = 'completed';
`
	rows, err := l.storage.Query(context.Background(), queryGetAllEventsRef)
	if err != nil {
		return nil, errors.Wrap(err,"unable to get events from storage")
	}
	defer rows.Close()

	var eventIDs []*string
	for rows.Next() {
		var s *string
		err = rows.Scan(&s)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to scan event with external reference = %s", s)
		}
		eventIDs = append(eventIDs, s)
	}

	return eventIDs, nil
}
