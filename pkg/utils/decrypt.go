package utils

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
)

func DecryptBody(c echo.Context, hashKey string, v interface{}) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(c.Request().Header.Get("X-Hash")), []byte(hashKey+string(body))); err != nil {
		return err
	}

	return nil
}
