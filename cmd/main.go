package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brendontj/lol-stats/app"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"github.com/brendontj/lol-stats/util"
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
		baseURI := util.GetEnvVariable("BASE_URI")
		token := util.GetEnvVariable("API_TOKEN")
		baseURIFeed := util.GetEnvVariable("BASE_URI_FEED")

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
	baseURI := util.GetEnvVariable("BASE_URI")
	token := util.GetEnvVariable("API_TOKEN")
	baseURIFeed := util.GetEnvVariable("BASE_URI_FEED")

	application := app.NewLolSportsClient(baseURI, token, baseURIFeed)
	application.Start()
	runningEvents := make(map[string]bool)
	for {
		liveGames := application.GetLiveGames()
		for _, lg := range liveGames.Data.Schedule.Events {
			if lg.State == "inProgress" {
				for _, g := range lg.Match.Games {
					if g.State == "inProgress" {
						_, ok := runningEvents[g.ID]
						if !ok {
							runningEvents[g.ID] = true
							for {
								currentGameData := application.GetCurrentLiveGame(g.ID)
								if currentGameData != nil {
									if len(currentGameData.Frames) > 0 && len(currentGameData.Frames[0].Participants) > 0{
										if currentGameData.Frames[0].Participants[0].TotalGoldEarned > 0 {
											break
										}
									}
								}
								time.Sleep(15 * time.Second)
							}

							for i := 0; i < 6; i++ {
								go func(m map[string]bool, lg lolsports.Events, gameMoment int) {
									time.Sleep(time.Duration(gameMoment*5) * time.Minute)
									type DataToBeSent struct {
										ID               string `json:"id"`
										TeamAWins        int    `json:"team_a_wins"`
										TeamALosses      int    `json:"team_a_losses"`
										TeamBWins        int    `json:"team_b_wins"`
										TeamBLosses      int    `json:"team_b_losses"`
										GameMoment       int    `json:"game_moment"`
										BlueSideTeamName string `json:"blue_side_team_name"`
										RedSideTeamName  string `json:"red_side_team_name"`
									}
									blueTeamName := ""
									redTeamName := ""
									if lg.Match.Teams[0].Side == "blue" {
										blueTeamName = lg.Match.Teams[0].Name
										redTeamName = lg.Match.Teams[1].Name
									} else {
										blueTeamName = lg.Match.Teams[1].Name
										redTeamName = lg.Match.Teams[0].Name
									}
									dataToBeSent := DataToBeSent{
										ID:               lg.ID,
										TeamAWins:        lg.Match.Teams[0].Record.Wins,
										TeamALosses:      lg.Match.Teams[0].Record.Losses,
										TeamBWins:        lg.Match.Teams[1].Record.Wins,
										TeamBLosses:      lg.Match.Teams[1].Record.Losses,
										GameMoment:       gameMoment,
										BlueSideTeamName: blueTeamName,
										RedSideTeamName:  redTeamName,
									}
									data, err := json.Marshal(dataToBeSent)
									if err != nil {
										panic(err)
									}
									resp, err := http.Post(util.GetEnvVariable("CONSUMER_API_ENDPOINT"), "application/json", bytes.NewBuffer(data))
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
	}
}
