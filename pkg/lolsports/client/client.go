package httpClient

import (
	"encoding/json"
	"fmt"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"log"
	"net/http"
	"time"
)

type EsportsAPIScrapper interface {
	GetLeagues(region string) (*lolsports.LeagueData, error)
	GetSchedule(region, leagueID, pageToken string) (*lolsports.ScheduleData, error)
	GetEventsLive(region string) (*lolsports.EventsLiveData, error)
	GetEventDetail(region, eventID string) (*lolsports.EventDetailData, error)
}

type PersistedDataClient struct {
	baseURI string
	token string
	httpClient *http.Client
}

func (h PersistedDataClient) GetLeagues(region string) (*lolsports.LeagueData, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprint(h.baseURI,"getLeagues?hl=", region),
		nil,
	)
	if err != nil {
		log.Printf("error creating HTTP request for get leagues of region %s: %v\n",region, err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-api-key", h.token)

	res, err := h.httpClient.Do(req)
	if err != nil {
		log.Printf("error sending HTTP request for get league of region %s: %v\n",region, err)
		return nil, err
	}

	leagueData := new(lolsports.LeagueData)

	if err := json.NewDecoder(res.Body).Decode(&leagueData); err != nil {
		log.Printf("error deserializing weather data\n")
		return nil, err
	}
	return leagueData, nil
}

func (h PersistedDataClient) GetSchedule(region, leagueID, pageToken string) (*lolsports.ScheduleData, error) {
	panic("implement me")
}

func (h PersistedDataClient) GetEventsLive(region string) (*lolsports.EventsLiveData, error) {
	panic("implement me")
}

func (h PersistedDataClient) GetEventDetail(region, eventID string) (*lolsports.EventDetailData, error) {
	panic("implement me")
}

func NewLolStatsClient(baseURI, token string) EsportsAPIScrapper {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &PersistedDataClient{
		baseURI: baseURI,
		token:   token,
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		},
	}
}