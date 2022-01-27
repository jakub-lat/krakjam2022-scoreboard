package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"krakjam2022_scoreboard/pkg/database"
	"krakjam2022_scoreboard/pkg/utils"
)

type Rest struct {
	db        *database.DB
	e         *echo.Echo
	secretKey string
}

func New(db *database.DB, secretKey string) *Rest {
	r := &Rest{db: db, e: echo.New(), secretKey: secretKey}

	r.e.Use(middleware.AddTrailingSlash())

	r.e.HTTPErrorHandler = func(err error, c echo.Context) {
		fmt.Println(err)
		_ = c.JSON(500, "sth went wrong")
	}

	r.e.GET("", r.GetTopScores)
	r.e.POST("/register", r.Register)
	r.e.GET("/player/:id", r.GetPlayerStats)
	r.e.GET("/player", r.GetCurrPlayerStats)
	r.e.POST("/run", r.PostRun)
	r.e.POST("/level", r.PostLevel)

	return r
}

func (r *Rest) Run(addr string) error {
	return r.e.Start(addr)
}

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

	run.Deaths += body.Deaths
	run.Kills += body.Kills
	run.Headshots += body.Headshots

	if err := r.db.Save(run).Error; err != nil {
		return err
	}

	return c.NoContent(200)
}

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
	return c.JSON(200, res)
}
