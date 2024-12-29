package utils

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func StringToPGUUID(uuidStr string) (pgtype.UUID, error) {
	var pgUUID pgtype.UUID

	// Parse the string into a UUID object
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return pgUUID, err
	}

	// Copy the bytes of the parsed UUID into the pgtype.UUID
	copy(pgUUID.Bytes[:], parsedUUID[:])

	return pgUUID, nil
}

func GetTempPath(root string, id string, filename string) string {
	return root + "/" + id + filename
}
