package app

import (
	"context"
	"fmt"
	"github.com/brendontj/lol-stats/pkg/lol-transformer/services"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type DataWorker interface {
	Start()
	TransformData()
	ExtractData()
	Close()
}

type dataWorker struct {
	dbPool                *pgxpool.Pool
	LolTransformerService services.Service
}

func NewDataWorker() DataWorker {
	return &dataWorker{dbPool: nil, LolTransformerService: nil}
}

func (d *dataWorker) Start() {
	dbPool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/lolstats?sslmode=disable&timezone=UTC") //Todo Add env vars
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializating the application: unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	d.LolTransformerService = services.NewLolService(dbPool)
	d.dbPool = dbPool
}

func (d *dataWorker) TransformData() {
	panic("implement me")
}

func (d *dataWorker) ExtractData() {
	if err := d.LolTransformerService.Extract("/Users/brendon/Projects/lol-stats/","stats"); err != nil {
		fmt.Println(err)
	}
}

func (d *dataWorker) Close() {
	d.dbPool.Close()
}