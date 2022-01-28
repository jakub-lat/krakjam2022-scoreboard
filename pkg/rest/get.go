package rest

import (
	"github.com/labstack/echo/v4"
	"krakjam2022_scoreboard/pkg/database"
	"krakjam2022_scoreboard/pkg/utils"
	"strconv"
)

func (r *Rest) GetPlayerStats(c echo.Context) error {
	id := c.Param("id")
	p := &database.Player{}
	if err := r.db.Preload("GameRuns").Preload("GameRuns.Levels").First(p, id).Error; err != nil {
		return err
	}

	return c.JSON(200, p)
}

func (r *Rest) GetCurrPlayerStats(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	p := &database.Player{}
	if err := r.db.Preload("GameRuns").Preload("GameRuns.Levels").Where("token = ?", token).First(p).Error; err != nil {
		return err
	}

	return c.JSON(200, p)
}

func (r *Rest) GetTopScores(c echo.Context) error {
	var res []database.Player
	err := r.db.Preload("GameRuns").Preload("GameRuns.Levels").Model(&database.Player{}).Find(&res).Error
	if err != nil {
		return err
	}

	data := struct {
		List []database.Player `json:"list"`
	}{res}
	return c.JSON(200, data)
}

func (r *Rest) GetTopScoresForLevel(c echo.Context) error {
	var res []database.GameRunLevel

	p, err := utils.Auth(r.db, c)
	if err != nil {
		return err
	}

	id := c.Param("id")
	limitStr := c.Param("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return err
	}

	err = r.db.Model(&database.GameRunLevel{}).
		Preload("Player").
		Distinct("player_id").
		Order("player_id, score desc, id").
		Where("level_id = ?", id).
		Limit(limit).
		Find(&res).Error
	if err != nil {
		return err
	}

	player := &database.GameRunLevel{}
	err = r.db.Preload("Player").Order("score desc").Where("player_id = ?", p.ID).First(player).Error
	if err != nil {
		return err
	}

	data := struct {
		Others []database.GameRunLevel `json:"others"`
		Player *database.GameRunLevel  `json:"player"`
	}{res, player}
	return c.JSON(200, data)
}
