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

	fmt.Println("Successfully inserted all leagues into database")
	return nil
}