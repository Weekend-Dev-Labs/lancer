package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetValidatedPayload(c echo.Context, payload interface{}) error {
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid payload")
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid payload")
	}
	return nil
}
