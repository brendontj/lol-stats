package lolsports

import "time"

type EventsLiveData struct {
	Data EventsLiveContent `json:"data"`
}

type EventsLiveResult struct {
	Outcome  string `json:"outcome"`
	GameWins int         `json:"gameWins"`
}

type EventsLiveRecord struct {
	Wins   int `json:"wins"`
	Losses int `json:"losses"`
}

type EventsLiveTeams struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Code   string `json:"code"`
	Image  string `json:"image"`
	Result EventsLiveResult `json:"result"`
	Record EventsLiveRecord `json:"record"`
}

type Events struct {
	ID        string    `json:"id,omitempty"`
	StartTime time.Time `json:"startTime"`
	State     string    `json:"state"`
	Type      string    `json:"type"`
	BlockName string    `json:"blockName"`
	League    League    `json:"league"`
	Match     Match     `json:"match"`
}

type EventsLiveContent struct {
	Schedule Schedule `json:"schedule"`
}