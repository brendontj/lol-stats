package main

import (
	"context"
	"fmt"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type application struct {
	dbPool           *pgxpool.Pool
	esportsApiClient lolsports.EsportsAPIScrapper
	feedApiClient    lolsports.FeedAPIScrapper
	lolService       lolsports.Service
}

func (a *application) init() {
	dbPool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/lolstats?sslmode=disable&timezone=UTC") //Todo Add env vars
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializating the application: unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	baseURI := "https://esports-api.lolesports.com/persisted/gw/" //Todo Add env vars
	token := "0TvQnueqKa5mxJntVWt0w4LpLfEkrV1Ta8rQBb9Z" //Todo Add env vars
	baseURIFeed := "https://feed.lolesports.com/livestats/v1/" //Todo Add env vars

	scrapperEsports := lolsports.NewLolStatsClient(baseURI, token)
	scrapperFeed := lolsports.NewLolFeedClient(baseURIFeed,token)

	a.lolService = lolsports.NewLolService(dbPool, scrapperEsports, scrapperFeed)

	a.dbPool = dbPool
	a.esportsApiClient = scrapperEsports
	a.feedApiClient = scrapperFeed
}

func (a *application) close() {
	a.feedApiClient.Close()
	a.esportsApiClient.Close()
	a.dbPool.Close()
}

func main() {
	app := application{}
	app.init()

	err := app.lolService.PopulateLeagues()
	if err != nil {
		panic("Unable to populate database with league of legends professional leagues")
	}

	app.close()
}