package lolsports

type LeagueData struct {
	ScheduleContent LeagueContent `json:"data"`
}

type League struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Region   string `json:"region,omitempty"`
	Image    string `json:"image"`
	Priority int    `json:"priority,omitempty"`
}

type LeagueContent struct {
	Leagues []League `json:"leagues"`
}