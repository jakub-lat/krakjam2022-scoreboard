package utils

import (
	"github.com/labstack/echo/v4"
	"krakjam2022_scoreboard/pkg/database"
)

func Auth(db *database.DB, c echo.Context) (*database.Player, error) {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return nil, echo.NewHTTPError(401, "Missing Authorization header")
	}

	p := &database.Player{}
	if err := db.Model(&database.Player{}).Where("token = ?", token).First(p).Error; err != nil {
		return nil, err
	}

	return p, nil
}
