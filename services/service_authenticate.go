package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type authInfo struct {
	ID string `json:"id"`
}

func (s *Services) authenticateFromServer(authToken string) (*authInfo, error) {
	authEndpoint := s.cfg.AuthEndpoint

	payload := map[string]string{
		"token": authToken,
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, authEndpoint, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Printf("[Lancer Error] Failed to create Requests: %s\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("[Lancer Error] Error making HTTP Requests: %s\n", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, err
	}

	var response authInfo

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		fmt.Printf("[Lancer Error] Server responds with unintended data content (%v)", err)

		return nil, err
	}

	return &response, nil
}
