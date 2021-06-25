package services

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
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