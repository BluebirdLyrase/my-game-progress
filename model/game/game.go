package game

import (
	"my-game-progress/model/environment"
)

type Game struct {
	year      int       `json:"year"`
	title     string    `json:"title"`
	genre     string    `json:"genre"`
	Remark    string    `json:"remark"`
	GameImage GameImage `json:"gameImage"`
}

type GameImage struct {
	Icon  string `json:"icon"`
	Cover string `json:"cover"`
}

type Playthrough struct {
	Difficulty    string   `json:"difficulty"`
	DateFinished  string   `json:"dateFinished"` //* DD/MM/YYYY
	Remark        string   `json:"remark"`
	Screenshots   []string `json:"screenshots"`
	Environment   environment.Environment
	GrandStrategy GrandStrategy
	RPG           RPG
}

// *
type Arcade struct {
	Score string `json:"Score"`
	Stat  []Stat `json:"stat"`
}

// *Strategy
type GrandStrategy struct {
	Campaign         string             `json:"campaign"`
	Faction          string             `json:"faction"`
	VictoryCondition []VictoryCondition `json:"victoryCondition"`
	Stat             []Stat             `json:"stat"` //*turn, Settlement Held, Battles Fought
}

type VictoryCondition struct {
	Title      string   `json:"title"`
	Objectives []string `json:"Objectives"`
	Success    bool     `json:"success"`
}

// * RPG
type RPG struct {
	Persona Persona `json:"persona"`
	Loadout Loadout `json:"loadout"`
	ARPG    ARPG    `json:"arpg"`
}

// * SMT
type Persona struct {
	Romance     string  `json:"romance"`
	SocialLinks []Stat  `json:"socialLinks"`
	SocialStats []Stat  `json:"socialStat"`
	PersonaList []Party `json:"personaList"`
	PartyList   []Party `json:"friendList"`
}

type Loadout struct {
	LoadoutSets []Equipment `json:"loadoutSets"`
}

type ARPG struct {
	Level      int         `json:"level"`
	Stat       []Stat      `json:"stat"`
	Skills     []Stat      `json:"skills"`
	Equipments []Equipment `json:"equipment"`
}

// * RPG Common
type Stat struct { //* any think like Level, Int, Str, Ma, LU, etc.
	Name  string `json:"name"`
	Image string `json:"image"`
	Value int    `json:"value"`
}

type Equipment struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Slot  string `json:"slot"`
}

type Party struct {
	Name       string      `json:"name"`
	Image      string      `json:"image"`
	Level      int         `json:"level"`
	Stat       []Stat      `json:"stat"`
	Skills     []Stat      `json:"skills"`
	Equipments []Equipment `json:"equipment"`
}
