package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"krakjam2022_scoreboard/pkg/database"
	"os"
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
	r.e.GET("/level/:id", r.GetTopScoresForLevel)
	r.e.GET("/clearalldata", r.ClearAllData)

	return r
}

func (r *Rest) Run(addr string) error {
	return r.e.Start(addr)
}

func (r *Rest) ClearAllData(c echo.Context) error {
	if c.QueryParam("asdf") != os.Getenv("CLEAR_ALL_DATA_KEY") {
		return echo.NewHTTPError(401, "Unauthorized")
	}

	r.db.Where("1 = 1").Delete(&database.Player{})
	r.db.Where("1 = 1").Delete(&database.GameRun{})
	r.db.Where("1 = 1").Delete(&database.GameRunLevel{})
	return c.JSON(200, "ok")
}
