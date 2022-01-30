package rest

import (
	"github.com/labstack/echo/v4"
	"krakjam2022_scoreboard/pkg/database"
	"krakjam2022_scoreboard/pkg/utils"
	"os"
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

func (r *Rest) GetAll(c echo.Context) error {
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

func (r *Rest) GetRun(c echo.Context) error {
	id := c.Param("id")
	run := &database.GameRun{}
	if err := r.db.Preload("Levels").First(run, id).Error; err != nil {
		return err
	}

	for _, x := range run.Levels {
		run.Kills += x.Kills
		run.Deaths += x.Deaths
		run.Score += x.Score
		run.Headshots += x.Headshots
	}

	return c.JSON(200, run)
}

func (r *Rest) GetTopScoresForLevel(c echo.Context) error {
	var res []database.GameRunLevel

	p, err := utils.Auth(r.db, c)
	if err != nil {
		return err
	}

	id := c.Param("id")
	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return err
	}

	err = r.db.Preload("Player").Raw(`SELECT DISTINCT ON ("player_id") *, ROW_NUMBER () OVER (ORDER BY score desc) AS position FROM "game_run_levels" WHERE level = ? AND "game_run_levels"."deleted_at" IS NULL ORDER BY player_id, score desc, id LIMIT ?`, id, limit).
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

func (r *Rest) GetTop(c echo.Context) error {

	p, err := utils.Auth(r.db, c)
	if err != nil {
		return err
	}

	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return err
	}

	var res []database.GameRun
	err = r.db.Preload("Player").Raw(`SELECT DISTINCT ON ("player_id") *, ROW_NUMBER () OVER (ORDER BY score desc) AS position FROM "game_runs" WHERE level = ? AND "game_runs"."deleted_at" IS NULL ORDER BY player_id, score desc, id LIMIT ?`, os.Getenv("MAX_LEVELS"), limit).
		Find(&res).Error
	if err != nil {
		return err
	}

	player := &database.GameRun{}
	err = r.db.Preload("Player").Order("score desc, id").Where("player_id = ?", p.ID).First(player).Error
	if err != nil {
		return err
	}

	data := struct {
		Others []database.GameRun `json:"others"`
		Player *database.GameRun  `json:"player"`
	}{res, player}
	return c.JSON(200, data)
}
