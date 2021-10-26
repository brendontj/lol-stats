package services

import (
	"context"
	lol_transformer "github.com/brendontj/lol-stats/pkg/lol-transformer"
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

func (s *service) fillPastStats(tx Transaction,gameID uuid.UUID, teamName string, numberOfPastGames int, gameMoment int, gameTime time.Time, teamOrder string) error {
	lastMatchStats, err := s.DB.GetLastMatchStats(teamName, numberOfPastGames, gameTime, gameMoment)
	if err != nil {
		return err
	}

	numberOfBaronsMean := 0
	numberOfDragonsMean := 0
	numberOfInhibitorsMean := 0
	numberOfTotalGoldMean := 0
	numberOfKillsMean := 0
	numberOfTowersMean := 0

	for _, lms := range lastMatchStats {
		if lms.TeamAName == teamName {
			numberOfBaronsMean += lms.TeamBlueTotalBarons
			numberOfDragonsMean += len(lms.TeamBlueDragons)
			numberOfInhibitorsMean += lms.TeamBlueTotalInhibitors
			numberOfTotalGoldMean += lms.TeamBlueTotalGold
			numberOfKillsMean += lms.TeamBlueTotalKills
			numberOfTowersMean += lms.TeamBlueTotalTowers
		} else if lms.TeamBName == teamName {
			numberOfBaronsMean += lms.TeamRedTotalBarons
			numberOfDragonsMean += len(lms.TeamRedDragons)
			numberOfInhibitorsMean += lms.TeamRedTotalInhibitors
			numberOfTotalGoldMean += lms.TeamRedTotalGold
			numberOfKillsMean += lms.TeamRedTotalKills
			numberOfTowersMean += lms.TeamRedTotalTowers
		}
	}

	statsInfo := lol_transformer.StatsInfo{
		NumberOfBaronsMean:     float64(numberOfBaronsMean)/float64(len(lastMatchStats)),
		NumberOfDragonsMean:    float64(numberOfDragonsMean)/float64(len(lastMatchStats)),
		NumberOfInhibitorsMean: float64(numberOfInhibitorsMean)/float64(len(lastMatchStats)),
		NumberOfTotalGoldMean:  float64(numberOfTotalGoldMean)/float64(len(lastMatchStats)),
		NumberOfKillsMean:      float64(numberOfKillsMean)/float64(len(lastMatchStats)),
		NumberOfTowersMean:     float64(numberOfTowersMean)/float64(len(lastMatchStats)),
	}

	return s.DB.UpdateMatchWithPastStats(tx,gameID,teamOrder, statsInfo, numberOfPastGames, gameMoment)
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
