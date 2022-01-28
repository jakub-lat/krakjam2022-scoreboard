package rest

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"krakjam2022_scoreboard/pkg/database"
	"krakjam2022_scoreboard/pkg/utils"
)

type registerReq struct {
	Name string `json:"name"`
}

func (r *Rest) Register(c echo.Context) error {
	var req registerReq
	if err := c.Bind(&req); err != nil {
		return err
	}
	token := utils.GenerateToken(20)

	p := &database.Player{
		Name:  req.Name,
		Token: token,
	}
	if err := r.db.Create(p).Error; err != nil {
		return err
	}

	return c.JSON(200, struct {
		*database.Player
		Token string `json:"token"`
	}{
		Player: p,
		Token:  p.Token,
	})
}

func (r *Rest) PostRun(c echo.Context) error {
	body := &database.GameRun{}
	err := utils.DecryptBody(c, r.secretKey, body)
	if err != nil {
		return err
	}

	p, err := utils.Auth(r.db, c)
	if err != nil {
		return err
	}

	body.PlayerID = p.ID

	err = r.db.Save(body).Error
	if err != nil {
		return err
	}

	return c.JSON(200, body)
}

func (r *Rest) PostLevel(c echo.Context) error {
	body := &database.GameRunLevel{}
	err := utils.DecryptBody(c, r.secretKey, body)
	if err != nil {
		return err
	}

	run := &database.GameRun{}
	if err := r.db.Model(&database.GameRun{}).Where("id = ?", body.GameRunID).First(run).Error; err != nil {
		return err
	}

	p, err := utils.Auth(r.db, c)
	if err != nil {
		return err
	}

	if run.PlayerID != p.ID {
		return echo.NewHTTPError(401, "Unauthorized")
	}

	body.PlayerID = p.ID
	body.Player = nil

	existingLevel := &database.GameRunLevel{}
	err = r.db.Model(&database.GameRunLevel{}).Where("level = ?", body.Level).First(existingLevel).Error
	if err == gorm.ErrRecordNotFound {
		if err := r.db.Create(body).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		body.ID = existingLevel.ID
		if err := r.db.Save(body).Error; err != nil {
			return err
		}
	}

	body.Player = p

	run.Deaths += body.Deaths
	run.Kills += body.Kills
	run.Headshots += body.Headshots

	if err := r.db.Save(run).Error; err != nil {
		return err
	}

	return c.JSON(200, body)
}
