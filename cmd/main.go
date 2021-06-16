package main

import (
	"context"
	"fmt"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"sync"

	//"sync"
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
		panic(err.Error())
	}
	fmt.Println("Successfully inserted all leagues into database")

	leagues, err := app.lolService.GetLeagues()
	if err != nil {
		panic("Unable to get leagues from database")
	}

	var wg sync.WaitGroup
	for _, league := range leagues {
		wg.Add(1)
		go func(leagueID string) {
			fmt.Printf("Starting inserts of games from league: %v \n", leagueID)
			defer wg.Done()
			err := app.lolService.PopulateDBScheduleOfLeague(leagueID)
			if err != nil {
				panic(err.Error())
			}
		}(league.ID)
	}
	wg.Wait()
	fmt.Println("Successfully inserted all events of schedule into database")

	eventIDs, err := app.lolService.GetEventsExternalRef()
	if err != nil {
		panic("Unable to get events from database")
	}

	for _, eventID := range eventIDs {
		wg.Add(1)
		go func(eventID string) {
			fmt.Printf("Inserting detailing of event: %v \n", eventID)
			defer wg.Done()
			err := app.lolService.PopulateDBWithEventDetail(eventID)
			if err != nil {
				panic(err.Error())
			}
		}(*eventID)
	}
	wg.Wait()
	fmt.Println("Successfully inserted all details of events into database")

	games, err := app.lolService.GetGamesReference()
	if err != nil {
		panic("Unable to get games from database")
	}

	for _, game := range games {
		wg.Add(1)
		go func(gameRef string) {
			fmt.Printf("Inserting game data of gameRef: %v \n", gameRef)
			defer wg.Done()
			err := app.lolService.PopulateDBWithGameData(gameRef)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(game)
	}
	wg.Wait()
	fmt.Println("Successfully inserted game data of all games into database")

	app.close()
}