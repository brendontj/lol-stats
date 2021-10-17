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
	"time"
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
					for i := 0; i < 6; i++ {
						go func(m map[string]bool, lg lolsports.Events, gameMoment int) {
							time.Sleep(time.Duration(gameMoment * 5) * time.Minute)
							type DataToBeSent struct {
								ID          string `json:"id"`
								TeamAWins   int `json:"team_a_wins"`
								TeamALosses int `json:"team_a_losses"`
								TeamBWins   int `json:"team_b_wins"`
								TeamBLosses int `json:"team_b_losses"`
								GameMoment  int `json:"game_moment"`
							}
							dataToBeSent := DataToBeSent{
								ID:          lg.ID,
								TeamAWins:   lg.Match.Teams[0].Record.Wins,
								TeamALosses: lg.Match.Teams[0].Record.Losses,
								TeamBWins:   lg.Match.Teams[1].Record.Wins,
								TeamBLosses: lg.Match.Teams[1].Record.Losses,
								GameMoment: gameMoment,
							}
							data, err := json.Marshal(dataToBeSent)
							if err != nil {
								panic(err)
							}
							resp, err := http.Post("http://localhost:8070/send_event_id/", "application/json", bytes.NewBuffer(data))
							if err != nil {
								panic(err)
							}
							if resp.StatusCode != http.StatusOK {
								fmt.Println(fmt.Sprintf("Error sending ID (%v)", lg.ID))
							}
						}(runningEvents, lg, i)
					}
				}
			}
		}
	}
}
