package lolsports

type ScheduleData struct {
	Data ScheduleContent `json:"data"`
}

type Pages struct {
	Older string `json:"older"`
	Newer string `json:"newer"`
}

type ScheduleLeague struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Result struct {
	Outcome  string `json:"outcome,omitempty"`
	GameWins int    `json:"gameWins"`
}

type Record struct {
	Wins   int `json:"wins"`
	Losses int `json:"losses"`
}

type Teams struct {
	TeamID     string `json:"id,omitempty"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Image  string `json:"image"`
	Result Result `json:"result"`
	Record Record `json:"record,omitempty"`
	Side string `json:"side,omitempty"`
}

type Strategy struct {
	Type  string `json:"type,omitempty"`
	Count int    `json:"count"`
}

type Match struct {
	ID       string        `json:"id,omitempty"`
	Flags    []string `json:"flags,omitempty"`
	Teams    []Teams       `json:"teams"`
	Strategy Strategy      `json:"strategy"`
	Games    []Games  `json:"games,omitempty"`
}

type Schedule struct {
	Pages  Pages    `json:"pages,omitempty"`
	Events []Events `json:"events"`
}

type ScheduleContent struct {
	Schedule Schedule `json:"schedule"`
}
