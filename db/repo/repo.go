package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/weekend-dev-labs/lancer/cache"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/types"
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

func (r *Repo) getTempPath(id string, filename string) string {
	return r.config.Store.Local.Temp + "/" + id + filename
}

func (r *Repo) CreateSession(session *types.SessionInfo) (*SessionCreateAck, error) {
	sessionKey := uuid.New().String()

	tempPath := r.getTempPath(sessionKey, session.FileName)

	sessionWithPath := &types.SessionInfo{
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

		return &SessionCreateAck{
			ID:       sessionKey,
			TempPath: tempPath,
		}, nil
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

	return &SessionCreateAck{
		ID:       ack.ID.String(),
		TempPath: tempPath,
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
