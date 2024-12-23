package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/weekend-dev-labs/lancer/types"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(url string) *Cache {
	if url != "" {
		opts, err := redis.ParseURL(url)

		if err != nil {
			log.Fatalf("[LANCER ERROR] Invalid Redis Connection URL")
		}

		client := redis.NewClient(opts)

		return &Cache{
			Client: client,
		}
	}

	return nil
}

func (c *Cache) CreateSession(session_id string, sessionInfo *types.SessionInfo) (bool, error) {
	key := "session" + ":" + session_id

	userJson, err := json.Marshal(sessionInfo)

	if err != nil {
		return false, fmt.Errorf("invalid session info")
	}

	err = c.Client.Set(context.Background(), key, userJson, time.Minute*60).Err()

	if err != nil {
		return false, fmt.Errorf("failed to set session info")
	}

	return true, nil
}

func (c *Cache) GetSessionInfo(session_id string) (*types.SessionInfo, error) {
	key := "session" + ":" + session_id

	val, err := c.Client.Get(context.Background(), key).Result()

	if err != nil {
		return nil, err
	}

	var sessionInfo types.SessionInfo

	err = json.Unmarshal([]byte(val), &sessionInfo)

	if err != nil {
		return nil, err
	}

	return &sessionInfo, nil
}
