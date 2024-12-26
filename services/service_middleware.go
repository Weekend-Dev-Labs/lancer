package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
)

type LancerValidator struct {
	Validator *validator.Validate
}

func (s *Services) middlewareAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeder := c.Request().Header.Get("authorization")

		splitedHeader := strings.Split(authHeder, " ")

		fmt.Println(splitedHeader)

		if len(splitedHeader) != 2 {
			return c.JSON(http.StatusForbidden, map[string]string{
				"err": "invalid auth header or auth header not present",
			})
		}

		// serverInfo, err := s.authenticateFromServer(splitedHeader[1])

		// if err != nil {
		// 	return c.JSON(http.StatusForbidden, map[string]string{
		// 		"err": "invalid auth header",
		// 	})
		// }

		c.Set(string(types.ContextAuthInfo), &authInfo{
			ID: "helloo",
		})

		return next(c)
	}
}

func (s *Services) middlewareSessionAuthenticator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionHeader := c.Request().Header.Get("x-session-token")

		if sessionHeader == "" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "invalid session to upload file",
			})
		}

		sessionInfo, err := utils.GetSessionInfo(sessionHeader, "secret key")

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid session token or session expired",
			})
		}

		c.Set(string(types.ContextSessionInfo), sessionInfo)

		return next(c)
	}
}

func (lv *LancerValidator) Validate(i interface{}) error {
	if err := lv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return nil
}
