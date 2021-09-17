package services

import "github.com/jackc/pgx/v4/pgxpool"

type Storage struct {
	pool *pgxpool.Pool
}