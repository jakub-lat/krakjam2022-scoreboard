package database

import (
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model `json:"-"`
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Token      string    `json:"-"`
	GameRuns   []GameRun `json:"gameRuns"`
}

const (
	GameModeEasy = iota
	GameModeNormal
	GameModeHard
	GameModeGod
)

type GameRun struct {
	gorm.Model `json:"-"`
	ID         int64          `json:"id"`
	PlayerID   int64          `json:"playerID"`
	Mode       int            `json:"mode"`
	StartTime  int64          `json:"startTime"`
	EndTime    int64          `json:"endTime"`
	Kills      int            `json:"kills"`
	Headshots  int            `json:"headshots"`
	Deaths     int            `json:"deaths"`
	Score      int            `json:"score"`
	Levels     []GameRunLevel `json:"levels"`
}

type GameRunLevel struct {
	gorm.Model `json:"-"`
	ID         int64   `json:"id"`
	GameRunID  int64   `json:"gameRunID"`
	PlayerID   int64   `json:"playerID"`
	Level      int     `json:"level"`
	StartTime  int64   `json:"startTime"`
	EndTime    int64   `json:"endTime"`
	Kills      int     `json:"kills"`
	Headshots  int     `json:"headshots"`
	Deaths     int     `json:"deaths"`
	Score      int     `json:"score"`
	Player     *Player `json:"player"`
	Position   int     `json:"position"`
}
