package rest

import (
	"database/sql"
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

var levelScoreboardSql = `with q as (
	select *, ROW_NUMBER() OVER (order by score desc) as position 
		from (select distinct on (player_id) * from game_run_levels order by player_id, score desc) as x
		where level = @level
		order by score desc
)
select distinct *
from 
(
	(select * from q limit @limit)
	UNION ALL 
	(select * from q where player_id=@player_id order by score desc limit 1)
) as res order by score desc`

func (r *Rest) GetTopScoresForLevel(c echo.Context) error {
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

	var res []database.GameRunLevel
	err = r.db.Preload("Player").Raw(levelScoreboardSql, sql.Named("level", id), sql.Named("limit", limit), sql.Named("player_id", p.ID)).
		Find(&res).Error
	if err != nil {
		return err
	}

	var player *database.GameRunLevel
	for _, x := range res {
		if x.Player.ID == p.ID {
			player = &x
			break
		}
	}

	data := struct {
		Others []database.GameRunLevel `json:"others"`
		Player *database.GameRunLevel  `json:"player"`
	}{res, player}
	return c.JSON(200, data)
}

var scoreboardSql = `with q as (
	select *, ROW_NUMBER() OVER (order by score desc) as position 
		from (select distinct on (player_id) * from game_runs order by player_id, score desc) as x
		where level = @level
		order by score desc
)
select distinct *
from 
(
	(select * from q limit @limit)
	UNION ALL 
	(select * from q where player_id=@player_id order by score desc limit 1)
) as res order by score desc`

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
	err = r.db.Preload("Player").Raw(scoreboardSql, sql.Named("level", os.Getenv("MAX_LEVELS")), sql.Named("limit", limit), sql.Named("player_id", p.ID)).
		Find(&res).Error

	if err != nil {
		return err
	}

	var player *database.GameRun
	for _, x := range res {
		if x.Player.ID == p.ID {
			player = &x
			break
		}
	}

	data := struct {
		Others []database.GameRun `json:"others"`
		Player *database.GameRun  `json:"player"`
	}{res, player}
	return c.JSON(200, data)
}
