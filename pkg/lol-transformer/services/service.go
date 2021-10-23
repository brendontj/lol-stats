package services

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service interface {
}

type service struct {
	storage *pgxpool.Pool
	DB      Storage
}

func NewLolService(pgStorage *pgxpool.Pool) Service {
	return &service{
		storage: pgStorage,
		DB:      Storage{pool: pgStorage},
	}
}

func (s *service) FillPastGamesWithHistoricData() error {

}

func (s *service) fillFormRatio(teamID string, numberOfGames int, )