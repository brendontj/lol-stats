package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brendontj/lol-stats/app"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go initRoutes(&wg)
	wg.Add(1)
	go handleLiveGames(&wg)
	wg.Wait()
}

func initRoutes(wg *sync.WaitGroup) {
	defer wg.Done()
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		 c.Data(http.StatusOK, "text/html", []byte("pong"))
	})

	router.GET("/sync_data", func(c *gin.Context) {
		baseURI := "https://esports-api.lolesports.com/persisted/gw/" //Todo Add env vars
		token := "0TvQnueqKa5mxJntVWt0w4LpLfEkrV1Ta8rQBb9Z"           //Todo Add env vars
		baseURIFeed := "https://feed.lolesports.com/livestats/v1/"    //Todo Add env vars

		application := app.NewLolSportsClient(baseURI, token, baseURIFeed)
		application.Start()
		application.PopulateHistoricalData()
		application.Close()

		c.Data(http.StatusOK, "text/html", []byte("synced"))
	})

	router.GET("/transform_data", func(c *gin.Context) {
		application := app.NewDataWorker()
		application.Start()
		application.TransformData()
		application.Close()

		c.Data(http.StatusOK, "text/html", []byte("data transformed"))
	})


	if err := router.Run(); err != nil {
		panic(err)
	}
}

func handleLiveGames(wg *sync.WaitGroup) {
	defer wg.Done()
	baseURI := "https://esports-api.lolesports.com/persisted/gw/" //Todo Add env vars
	token := "0TvQnueqKa5mxJntVWt0w4LpLfEkrV1Ta8rQBb9Z"           //Todo Add env vars
	baseURIFeed := "https://feed.lolesports.com/livestats/v1/"    //Todo Add env vars

	application := app.NewLolSportsClient(baseURI, token, baseURIFeed)
	application.Start()
	runningEvents := make(map[string]bool)
	for {
		liveGames := application.GetLiveGames()
		for _, lg := range liveGames.Data.Schedule.Events {
			if lg.State == "inProgress" {
				_, ok := runningEvents[lg.ID]
				if !ok {
					runningEvents[lg.ID] = true

					go func(m map[string]bool, lg lolsports.Events) {
						type DataToBeSent struct {
							id string
							team_a_wins int
							team_a_losses int
							team_b_wins int
							team_b_losses int
						}

						data, _ := json.Marshal(DataToBeSent{
							id: lg.ID,
							team_a_wins: lg.Match.Teams[0].Record.Wins,
							team_a_losses: lg.Match.Teams[1].Record.Losses,
							team_b_wins: lg.Match.Teams[1].Record.Wins,
							team_b_losses: lg.Match.Teams[1].Record.Losses,
						})

						resp, err := http.Post("localhost:8070/send_event_id", "application/json", bytes.NewBuffer(data))
						if err != nil {
							panic(err)
						}
						if resp.StatusCode != http.StatusOK {
							fmt.Println(fmt.Sprintf("Error sending ID (%v)", lg.ID))
							delete(m,lg.ID)
						}
					}(runningEvents, lg)
				}
			}
		}
	}
}
