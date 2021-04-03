package httpClient

import (
	"encoding/json"
	"fmt"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"log"
	"net/http"
	"time"
)

const EmptyField = ""

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
	if leagueID != EmptyField {
		leagueID = fmt.Sprint("&leagueId=", leagueID)
	}

	if pageToken != EmptyField {
		pageToken = fmt.Sprint("&pageToken=", pageToken)
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprint(h.baseURI,"getSchedule?hl=", region, leagueID, pageToken),
		nil,
	)
	if err != nil {
		log.Printf("error creating HTTP request for get schedule of region %s: %v\n",region, err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-api-key", h.token)

	res, err := h.httpClient.Do(req)
	if err != nil {
		log.Printf("error sending HTTP request for get league of region %s: %v\n",region, err)
		return nil, err
	}

	scheduleData := new(lolsports.ScheduleData)

	if err := json.NewDecoder(res.Body).Decode(&scheduleData); err != nil {
		log.Printf("error deserializing weather data\n")
		return nil, err
	}
	return scheduleData, nil
}

func (h PersistedDataClient) GetEventsLive(region string) (*lolsports.EventsLiveData, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprint(h.baseURI,"getLive?hl=", region),
		nil,
	)
	if err != nil {
		log.Printf("error creating HTTP request for get events live of region %s: %v\n",region, err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-api-key", h.token)

	res, err := h.httpClient.Do(req)
	if err != nil {
		log.Printf("error sending HTTP request for get events live of region %s: %v\n",region, err)
		return nil, err
	}

	eventsLiveData := new(lolsports.EventsLiveData)

	if err := json.NewDecoder(res.Body).Decode(&eventsLiveData); err != nil {
		log.Printf("error deserializing weather data\n")
		return nil, err
	}
	return eventsLiveData, nil
}

func (h PersistedDataClient) GetEventDetail(region, eventID string) (*lolsports.EventDetailData, error) {
	if eventID != EmptyField {
		eventID = fmt.Sprint("&id=", eventID)
	}
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprint(h.baseURI,"getEventDetails?hl=", region, eventID),
		nil,
	)
	if err != nil {
		log.Printf("error creating HTTP request for get event details of region %s and event %s: %v\n",region, eventID, err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-api-key", h.token)

	res, err := h.httpClient.Do(req)
	if err != nil {
		log.Printf("error sending HTTP request for get event details of region %s and event %s: %v\n",region, eventID, err)
		return nil, err
	}

	eventDetailData := new(lolsports.EventDetailData)

	if err := json.NewDecoder(res.Body).Decode(&eventDetailData); err != nil {
		log.Printf("error deserializing weather data\n")
		return nil, err
	}
	return eventDetailData, nil
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