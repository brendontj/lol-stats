package app

import (
	"context"
	"fmt"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"github.com/brendontj/lol-stats/pkg/lolsports/services"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"sync"
)

type LolSportsClient interface {
	Start()
	PopulateHistoricalData()
	Close()
}

type lolSportsClient struct {
	LolService    services.Service
	dbPool           *pgxpool.Pool
	esportsApiClient lolsports.EsportsAPIScrapper
	feedApiClient lolsports.FeedAPIScrapper
}

func NewLolSportsClient(baseURI, token, baseURIFeed string) LolSportsClient {
	return &lolSportsClient{
		LolService:       nil,
		dbPool:           nil,
		esportsApiClient: lolsports.NewLolStatsClient(baseURI, token),
		feedApiClient:    lolsports.NewLolFeedClient(baseURIFeed, token),
	}
}

func (a *lolSportsClient) Start() {
	dbPool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/lolstats?sslmode=disable&timezone=UTC") //Todo Add env vars
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializating the application: unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	a.LolService = services.NewLolService(dbPool, a.esportsApiClient, a.feedApiClient)
	a.dbPool = dbPool
}

func (a *lolSportsClient) PopulateHistoricalData() {
	err := a.LolService.PopulateLeagues()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully inserted all leagues into database")

	leagues, err := a.LolService.GetLeagues()
	if err != nil {
		panic("Unable to get leagues from database")
	}

	var wg sync.WaitGroup
	for i, league := range leagues {
		wg.Add(1)
		go func(wg *sync.WaitGroup, leagueID string, i int) {
			fmt.Printf("[Worker %v] Starting inserts of games from league: %v \n",i, leagueID)
			defer wg.Done()
			err := a.LolService.PopulateDBScheduleOfLeague(leagueID)
			if err != nil {
				panic(err.Error())
			}
		}(&wg, league.ID, i)
	}
	fmt.Println("[Main process]: Waiting for workers to finish")
	wg.Wait()
	fmt.Println("[Main process]: Successfully inserted all events of schedule into database")

	eventIDs, err := a.LolService.GetEventsExternalRef()
	if err != nil {
		panic("Unable to get events from database")
	}

	for i, eventID := range eventIDs {
		wg.Add(1)
		go func(wg *sync.WaitGroup, eventID string, i int) {
			fmt.Printf("[Worker %v] Inserting detailing of event: %v \n", i, eventID)
			defer wg.Done()
			err := a.LolService.PopulateDBWithEventDetail(eventID)
			if err != nil {
				panic(err.Error())
			}
		}(&wg, *eventID, i)
	}
	fmt.Println("[Main process]: Waiting for workers to finish")
	wg.Wait()
	fmt.Println("[Main process]: Successfully inserted all details of events into database")

	games, err := a.LolService.GetGamesReference()
	if err != nil {
		panic("Unable to get games from database")
	}

	for i, game := range games {
		wg.Add(1)
		go func(wg *sync.WaitGroup, gameRef string, i int) {
			fmt.Printf("[Worker %v] Inserting game data of gameRef: %v \n", i, gameRef)
			defer wg.Done()
			if err := a.LolService.PopulateDBWithGameData(gameRef); err != nil {
				panic(err.Error())
			}
		}(&wg, game, i)
	}
	fmt.Println("[Main process]: Waiting for workers to finish")
	wg.Wait()
	fmt.Println("[Main process]: Successfully inserted game data of all games into database")
}

func (a *lolSportsClient) Close() {
	a.feedApiClient.Close()
	a.esportsApiClient.Close()
	a.dbPool.Close()
}