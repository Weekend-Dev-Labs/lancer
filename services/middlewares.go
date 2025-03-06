package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
)

type LancerValidator struct {
	Validator *validator.Validate
}

func (s *Services) authClientServerToken(c echo.Context) (echo.Context, error) {

	payload := new(types.CreateSessionPayload)

	if err := utils.GetValidatedPayload(c, payload); err != nil {
		return c, fmt.Errorf("invalid session create payload")
	}

	if s.cfg.ServerAuth {
		authHeder := c.Request().Header.Get("authorization")

		splitedHeader := strings.Split(authHeder, " ")

		if len(splitedHeader) != 2 {
			return c, fmt.Errorf("missing auth token")
		}

		jsonPayload, err := json.Marshal(payload)

		if err != nil {
			return c, fmt.Errorf("invalid session create payload")
		}

		req, err := http.NewRequest(http.MethodPost, s.cfg.AuthEndpoint, bytes.NewReader(jsonPayload))

		if err != nil {
			fmt.Println(err.Error())
			return c, fmt.Errorf("failed to make request to auth server")
		}

		req.Header.Set("authorization", "Bearer "+splitedHeader[1])
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{
			Timeout: 10 * time.Second,
		}

		res, err := client.Do(req)

		if err != nil {
			fmt.Println("Request Failed")
			return c, fmt.Errorf("failed to make request to auth server")
		}

		if res.StatusCode != 200 {
			return c, fmt.Errorf("invalid auth token")
		}

		c.Set(string(types.ContextAuthInfo), &authInfo{
			ID: "helloo",
		})

		c.Set(string(types.ContextSessionPayload), payload)

		return c, nil
	}

	fmt.Println("Without Authentication")

	c.Set(string(types.ContextAuthInfo), &authInfo{
		ID: "lancer",
	})

	c.Set(string(types.ContextSessionPayload), payload)

	return c, nil
}

func (s *Services) authSessionToken(c echo.Context) (echo.Context, bool) {
	sessionHeader := c.Request().Header.Get("x-session-token")

	// if sessionHeader == "" {
	// 	return c.JSON(http.StatusForbidden, map[string]string{
	// 		"error": "invalid session to upload file",
	// 	})
	// }

	if sessionHeader == "" {
		return c, false
	}

	sessionInfo, err := utils.GetSessionInfo(sessionHeader, "secret key")

	if err != nil {
		return c, false
	}

	c.Set(string(types.ContextSessionInfo), sessionInfo)

	return c, true
}

func (s *Services) authWebToken(c echo.Context) (echo.Context, bool) {
	adminHeader := c.Request().Header.Get("Authorization")

	if adminHeader == "" {
		return c, false
	}

	splittedAdminHeader := strings.Split(adminHeader, " ")

	if len(splittedAdminHeader) < 2 {
		return c, false
	}

	adminInfo, err := utils.GetAdminInfo(splittedAdminHeader[1], s.cfg.AdminTokenSigningSecret)

	if err != nil {
		return c, false
	}

	c.Set(string(types.ContextWebToken), adminInfo)

	return c, true
}

func (s *Services) middlewareAuth(authType []types.AuthKeys, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var isTokenValid bool
		for _, auth := range authType {
			switch auth {
			case types.AuthClientServerToken:
				fmt.Println("Client Server Token")
				handlerC, err := s.authClientServerToken(c)

				if err != nil {
					return c.JSON(http.StatusForbidden, map[string]string{
						"err": err.Error(),
					})
				}
				isTokenValid = true

				return next(handlerC)

			case types.AuthSessionToken:
				handlerC, isTokenCorrect := s.authSessionToken(c)

				if isTokenCorrect {
					return next(handlerC)
				}

				isTokenValid = false

			case types.AuthWebToken:
				handlerC, isTokenCorrect := s.authWebToken(c)

				if isTokenCorrect {
					return next(handlerC)
				}

				isTokenValid = false
			}
		}

		if !isTokenValid {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "invalid tokens to access the endpoint",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "something went wrong",
		})
	}
}

func (lv *LancerValidator) Validate(i interface{}) error {
	if err := lv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return nil
}
