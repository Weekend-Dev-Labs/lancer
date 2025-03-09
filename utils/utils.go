package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"

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

func GetJsonStruct(byteArr []byte, val interface{}) error {
	return json.Unmarshal(byteArr, val)
}

func GenerateSecret(length int) (string, error) {
	// Create a byte slice to store the random data
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	// Convert the byte slice into a base64 encoded string (or use hex if preferred)
	encodedSecret := base64.URLEncoding.EncodeToString(secret)
	return encodedSecret, nil
}

func GetMimetypeByPath(path string) string {
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "text/plain" // Default to plain text if the MIME type is unknown
	}
	return mimeType
}