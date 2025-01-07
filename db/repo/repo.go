package repo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/weekend-dev-labs/lancer/cache"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/types"
	"golang.org/x/crypto/bcrypt"
)

type Repo struct {
	db         *db.Queries
	redisCache *cache.Cache
	config     *config.LancerConfig
}

type SessionCreateAck struct {
	ID       string
	TempPath string
}

func NewRepo(db *db.Queries, redisCache *cache.Cache, config *config.LancerConfig) *Repo {
	return &Repo{
		db:         db,
		redisCache: redisCache,
		config:     config,
	}
}

func CreateInitialUser(cfg *config.LancerConfig, query *db.Queries) {
	isUserExists, err := query.CheckEmailExists(context.TODO(), cfg.Auth.Email)

	if err != nil {
		log.Fatalf("[Lancer Error] Failed to check admin user (%v)", err)
	}

	if !isUserExists {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Auth.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Fatalf("[Lancer Error] Failed to hash password for admin user (%v)", err)
		}

		_, err = query.CreateUser(context.TODO(), db.CreateUserParams{
			Email:    cfg.Auth.Email,
			Password: string(hashedPassword),
		})

		if err != nil {
			log.Fatalf("[Lancer Error] Failed to create admin user")
		}
	}
}

func CreateOrGetInitialMetrics(cfg *config.LancerConfig, query *db.Queries) {
	metrics, err := query.GetFirstCreatedMetrics(context.Background())

	if err != nil {
		newMetrics, err := query.InsertMetrics(context.TODO(), db.InsertMetricsParams{})

		if err != nil {
			log.Fatalf("[Lancer Error] Failed to create metrics for the app (%v)", err)
		}

		cfg.MetricsID = newMetrics.ID
	}

	cfg.MetricsID = metrics.ID
}

func (r *Repo) getTempPath(id string, filename string) string {
	return r.config.Store.Local.Temp + "/" + id + filename
}

func (r *Repo) CreateSession(session *types.SessionInfo) (*types.SessionInfo, error) {
	sessionKey := uuid.New().String()

	tempPath := r.getTempPath(sessionKey, session.FileName)

	sessionWithPath := &types.SessionInfo{
		ID:           sessionKey,
		FileSize:     session.FileSize,
		ChunkSize:    session.ChunkSize,
		MaxChunk:     session.MaxChunk,
		FileName:     session.FileName,
		TempPath:     tempPath,
		OwnerID:      session.OwnerID,
		CurrentChunk: 0,
	}

	if r.config.Redis != "" {
		_, err := r.redisCache.CreateSession(sessionKey, sessionWithPath)

		if err != nil {
			return nil, err
		}

		return sessionWithPath, nil
	}

	dbSessionData := db.CreateSessionParams{
		FileSize:  session.FileSize,
		ChunkSize: session.ChunkSize,
		MaxChunk:  session.MaxChunk,
		FileName: pgtype.Text{
			String: session.FileName,
			Valid:  true,
		},
		OwnerID: pgtype.Text{
			String: session.OwnerID,
			Valid:  true,
		},
		Provider: pgtype.Text{
			String: string(session.Provider),
			Valid:  true,
		},
	}

	ack, err := r.db.CreateSession(context.TODO(), dbSessionData)

	if err != nil {
		return nil, err
	}

	tempPath = r.getTempPath(ack.ID.String(), ack.FileName.String)

	if err := r.db.UpdateSession(context.TODO(), db.UpdateSessionParams{
		FileName: ack.FileName,
		TempPath: pgtype.Text{
			String: r.getTempPath(ack.ID.String(), ack.FileName.String),
			Valid:  true,
		},
		ID: ack.ID,
	}); err != nil {
		return nil, nil
	}

	return &types.SessionInfo{
		ID:           ack.ID.String(),
		TempPath:     tempPath,
		FileSize:     ack.FileSize,
		ChunkSize:    ack.ChunkSize,
		MaxChunk:     ack.MaxChunk,
		FileName:     ack.FileName.String,
		OwnerID:      ack.OwnerID.String,
		CurrentChunk: ack.CurrentChunk.Int64,
		Provider:     types.UploaderProvider(ack.Provider.String),
	}, nil
}

func (r *Repo) DeleteSession(id string) error {
	if r.config.Redis != "" {
		r.redisCache.RemoveSession(id)
		return nil
	}

	sessionUuid, err := uuid.Parse(id)

	if err != nil {
		return err
	}

	if err := r.db.DeleteSession(context.TODO(), pgtype.UUID{
		Bytes: sessionUuid,
		Valid: true,
	}); err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetSessionById(id string) (*types.SessionInfo, error) {
	if r.config.Redis != "" {
		sessionInfo, err := r.redisCache.GetSessionInfo(id)

		if err != nil {
			return nil, err
		}

		return sessionInfo, nil
	}

	sessionId, _ := uuid.Parse(id)

	sessionInfo, err := r.db.FindSessionById(context.TODO(), pgtype.UUID{
		Bytes: sessionId,
		Valid: true,
	})

	if err != nil {
		return nil, err
	}

	return &types.SessionInfo{
		FileSize:     sessionInfo.FileSize,
		ChunkSize:    sessionInfo.ChunkSize,
		MaxChunk:     sessionInfo.MaxChunk,
		FileName:     sessionInfo.FileName.String,
		TempPath:     sessionInfo.TempPath.String,
		OwnerID:      sessionInfo.OwnerID.String,
		CurrentChunk: sessionInfo.CurrentChunk.Int64,
	}, nil
}

func (r *Repo) UpdateSessionById(id string, session *types.SessionInfo) error {
	if r.config.Redis != "" {
		_, err := r.redisCache.CreateSession(id, session)

		if err != nil {
			return err
		}

		return nil
	}

	sessionId, err := uuid.Parse(id)

	if err != nil {
		return err
	}

	dbSessionInfo := db.UpdateSessionParams{
		TempPath: pgtype.Text{
			String: session.TempPath,
			Valid:  true,
		},
		CurrentChunk: pgtype.Int8{
			Int64: session.CurrentChunk,
			Valid: true,
		},
		FileName: pgtype.Text{
			String: session.FileName,
			Valid:  true,
		},
		Provider: pgtype.Text{
			String: string(session.Provider),
			Valid:  true,
		},
		ID: pgtype.UUID{
			Bytes: sessionId,
			Valid: true,
		},
	}

	if err := r.db.UpdateSession(context.TODO(), dbSessionInfo); err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetSessions(params *db.PaginateSessionsParams) []*types.SessionInfo {
	if r.config.Redis != "" {
		var cursor uint64
		var keys []string
		var err error

		var sessionInfo []*types.SessionInfo

		for {
			keys, cursor, err = r.redisCache.Client.Scan(context.Background(), cursor, "session", int64(params.Limit)).Result()

			if err != nil {
				return nil
			}

			for _, key := range keys {
				fmt.Println(key)

				splitedKey := strings.Split(key, ":")

				if len(splitedKey) == 2 {
					sessionId := splitedKey[1]
					session, err := r.redisCache.GetSessionInfo(sessionId)

					if err != nil {
						return nil
					}

					sessionInfo = append(sessionInfo, session)
				}
			}

			if cursor == 0 {
				break
			}
		}

		return sessionInfo
	}

	var sessionInfo []*types.SessionInfo

	sessions, err := r.db.PaginateSessions(context.TODO(), db.PaginateSessionsParams{
		Limit:  100,
		Offset: 0,
	})

	if err != nil {
		return nil
	}

	for _, val := range sessions {
		sessionInfo = append(sessionInfo, &types.SessionInfo{
			ID:           val.ID.String(),
			FileSize:     val.FileSize,
			ChunkSize:    val.ChunkSize,
			MaxChunk:     val.MaxChunk,
			FileName:     val.FileName.String,
			TempPath:     val.TempPath.String,
			OwnerID:      val.OwnerID.String,
			CurrentChunk: val.CurrentChunk.Int64,
			Provider:     types.UploaderProvider(val.Provider.String),
		})
	}

	return sessionInfo
}
