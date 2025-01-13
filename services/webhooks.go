package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookEvent string

type Webhook struct {
	endpoint      string
	signingSecret string
}

const (
	EventSessionCreate     = WebhookEvent("SESSION_CREATED")
	EventSessionCancelled  = WebhookEvent("SESSION_CANCELLED")
	EventSessionCompleteed = WebhookEvent("SESSION_COMPLETED")

	EventFileUpload = WebhookEvent("FILE_UPLOAD")
	EventFileDelete = WebhookEvent("FILE_DELETE")
)

func NewWebhookNotifier(endpoint string, signingSecret string) *Webhook {
	return &Webhook{
		endpoint:      endpoint,
		signingSecret: signingSecret,
	}
}

func (wh *Webhook) SendEvent(event WebhookEvent, payload interface{}) error {
	if wh.endpoint != "" {
		data := map[string]interface{}{
			"event": event,
			"data":  payload,
		}

		jsonData, err := json.Marshal(data)

		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, wh.endpoint, bytes.NewBuffer(jsonData))

		if err != nil {
			return err
		}

		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-timestamp", timestamp)
		req.Header.Set("x-signature", getSignature(string(jsonData), timestamp, wh.signingSecret))

		client := http.Client{
			Timeout: 10 * time.Second,
		}

		_, err = client.Do(req)

		if err != nil {
			return err
		}

		// TODO : metrics for webhook success & failure rates

		return nil
	}

	return nil
}

func getSignature(payload, timestamp, secret string) string {
	message := timestamp + "." + payload
	mac := hmac.New(sha256.New, []byte(secret))

	mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil))
}
