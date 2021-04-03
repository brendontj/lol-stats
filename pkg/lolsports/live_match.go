package lolsports

import "time"

type LiveMatchData struct {
	Frames []Frames `json:"frames"`
}

type PerkMetadata struct {
	StyleID    int   `json:"styleId"`
	SubStyleID int   `json:"subStyleId"`
	Perks      []int `json:"perks"`
}

type Participants struct {
	ParticipantID       int          `json:"participantId"`
	Level               int          `json:"level"`
	Kills               int          `json:"kills"`
	Deaths              int          `json:"deaths"`
	Assists             int          `json:"assists"`
	TotalGoldEarned     int          `json:"totalGoldEarned"`
	CreepScore          int          `json:"creepScore"`
	KillParticipation   float64      `json:"killParticipation"`
	ChampionDamageShare float64      `json:"championDamageShare"`
	WardsPlaced         int          `json:"wardsPlaced"`
	WardsDestroyed      int          `json:"wardsDestroyed"`
	AttackDamage        int          `json:"attackDamage"`
	AbilityPower        int          `json:"abilityPower"`
	CriticalChance      float64      `json:"criticalChance"`
	AttackSpeed         int          `json:"attackSpeed"`
	LifeSteal           int          `json:"lifeSteal"`
	Armor               int          `json:"armor"`
	MagicResistance     int          `json:"magicResistance"`
	Tenacity            float64      `json:"tenacity"`
	Items               []int        `json:"items"`
	PerkMetadata        PerkMetadata `json:"perkMetadata"`
	Abilities           []string     `json:"abilities"`
	TotalGold           int          `json:"totalGold,omitempty"`
	CurrentHealth       int          `json:"currentHealth,omitempty"`
	MaxHealth           int          `json:"maxHealth,omitempty"`
}

type Frames struct {
	Rfc460Timestamp time.Time      `json:"rfc460Timestamp"`
	Participants    []Participants `json:"participants,omitempty"`
	GameState       string         `json:"gameState,omitempty"`
	BlueTeam        Team       `json:"blueTeam,omitempty"`
	RedTeam         Team        `json:"redTeam,omitempty"`
}

type Team struct {
	TotalGold    int            `json:"totalGold"`
	Inhibitors   int            `json:"inhibitors"`
	Towers       int            `json:"towers"`
	Barons       int            `json:"barons"`
	TotalKills   int            `json:"totalKills"`
	Dragons      []interface{}  `json:"dragons"`
	Participants []Participants `json:"participants"`
}