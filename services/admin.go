package services

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *Services) serviceLoginAdmin(c echo.Context) error {

	payload := new(types.AdminUserPayload)

	if err := utils.GetValidatedPayload(c, payload); err != nil {
		return err
	}

	isUserExists, err := s.db.GetUserByEmail(context.TODO(), payload.Email)

	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(isUserExists.Password), []byte(payload.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid password",
		})
	}

	token, err := utils.GetAdminToken(&utils.AdminClaims{
		AdminID: isUserExists.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}, s.cfg.AdminTokenSigningSecret)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (s *Services) registerAdminService() {
	group := s.e.Group("/admin")

	group.POST("/login", s.serviceLoginAdmin)
}
