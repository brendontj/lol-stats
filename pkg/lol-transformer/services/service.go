package services

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
	"time"
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
	matches, err := s.DB.GetAllMatchesWithoutProcessedData()
	if err != nil {
		return err
	}

	tx, err := s.storage.Begin(context.Background())
	if err != nil {
		return err
	}

	for _, m := range matches {
		if err := s.fillFormRatio(tx,m.ID, m.TeamAName, 5, m.GameTime, "A"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillFormRatio(tx,m.ID, m.TeamBName, 5, m.GameTime, "B"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillFormRatio(tx,m.ID, m.TeamAName, 3, m.GameTime, "A"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillFormRatio(tx,m.ID, m.TeamBName, 3, m.GameTime, "B"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillPastStats(tx,m.ID, m.TeamAName, 5, 4, m.GameTime, "A"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillPastStats(tx,m.ID, m.TeamBName, 5, 4, m.GameTime, "B"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillPastStats(tx,m.ID, m.TeamAName, 3, 4, m.GameTime, "A"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillPastStats(tx,m.ID, m.TeamBName, 3, 4, m.GameTime, "B"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillPastStats(tx,m.ID, m.TeamAName, 5, 6, m.GameTime, "A"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillPastStats(tx,m.ID, m.TeamBName, 5, 6, m.GameTime, "B"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillPastStats(tx,m.ID, m.TeamAName, 3, 6, m.GameTime, "A"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err := s.fillPastStats(tx,m.ID, m.TeamBName, 3, 6, m.GameTime, "B"); err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}
	}
	if err := tx.Commit(context.Background()); err != nil {
		_ = tx.Rollback(context.Background())
	}

	return nil
}

func (s *service) fillFormRatio(tx Transaction,gameID uuid.UUID, teamName string, numberOfPastGames int, gameTime time.Time, teamOrder string) error {
	lastMatchResults, err := s.DB.GetLastMatchResults(tx, teamName, numberOfPastGames, gameTime)
	if err != nil {
		return err
	}

	numberOfWins := 0
	for _, lmr := range lastMatchResults {
		if lmr.TeamAName == teamName {
			if isMatchWinner(lmr.BestOf, lmr.TeamAGameWins) {
				numberOfWins += 1
			}
		} else if lmr.TeamBName == teamName {
			if isMatchWinner(lmr.BestOf, lmr.TeamBGameWins) {
				numberOfWins += 1
			}
		}
	}
	ratio := float64(numberOfWins) / float64(len(lastMatchResults))

	return s.DB.UpdateMatchWithWinnerRatio(tx, gameID, teamOrder, ratio, numberOfPastGames)
}

func (s *service) fillPastStats(tx Transaction,gameID uuid.UUID, teamName string, numberOfPastGames int, gameMoment int, gameTime time.Time) error {
	lastMatchStats, err := s.DB.GetLastMatchStats(tx, teamName, numberOfPastGames, gameTime)
	if err != nil {
		return err
	}
}

func isMatchWinner(bestOf int, numberOfWins int) bool {
	switch bestOf {
	case 1:
		if numberOfWins == 1 {
			return true
		}
		break
	case 2, 3:
		if numberOfWins == 2 {
			return true
		}
		break
	case 5:
		if numberOfWins == 3 {
			return true
		}
		break
	default:
		return false
	}
	return false
}
