package types

import (
	"github.com/google/uuid"
	"github.com/weekend-dev-labs/lancer/db"
)

type ContextKeys string
type AuthKeys string
type UploaderProvider string

const (
	ContextAuthInfo       = ContextKeys("auth-info")
	ContextSessionInfo    = ContextKeys("session-info")
	ContextWebToken       = ContextKeys("webtoken-info")
	ContextSessionPayload = ContextKeys("session-payload")
)

const (
	AuthWebToken          = AuthKeys("web-token")
	AuthSessionToken      = AuthKeys("session-token")
	AuthClientServerToken = AuthKeys("client-token")
	AuthServerCredentials = AuthKeys("server-credentials")
)

const (
	UploaderAws   = UploaderProvider("AWS")
	UploaderLocal = UploaderProvider("LOCAL")
)

const (
	AppName       = "lancer"
	AppConfigFile = "lancer.yaml"
	AppSecrets    = "secrets"
	AppHistory    = "history"
)

type CreateSessionPayload struct {
	FileSize  int64            `json:"file_size" validate:"required"`
	ChunkSize int64            `json:"chunk_size" validate:"required"`
	MaxChunk  int64            `json:"max_chunk" validate:"required"`
	FileName  string           `json:"file_name" validate:"required"`
	Provider  UploaderProvider `json:"provider" validate:"required"`
}

type SessionTokenPayload struct {
	SessionToken string `json:"session_token" validate:"required"`
}

type SessionInfo struct {
	ID           string           `json:"id"`
	FileSize     int64            `json:"file_size"`
	ChunkSize    int64            `json:"chunk_size"`
	MaxChunk     int64            `json:"max_chunk"`
	FileName     string           `json:"file_name"`
	TempPath     string           `json:"temp_path"`
	OwnerID      string           `json:"owner_id"`
	CurrentChunk int64            `json:"current_chunk"`
	Provider     UploaderProvider `json:"provider"`
}

type UploaderChunkPayload struct {
	ChunkCount int64  `validate:"required"`
	Checksum   string `validate:"required"`
}

type AdminUserPayload struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type UploadDeletePayload struct {
	ID []uuid.UUID `json:"id" validate:"required"`
}

type UploadQueryInfo struct {
	Limit int64 `query:"size"`
	Page  int64 `query:"page"`
}

type DeleteUploadFileList struct {
	db.DeleteDocumentsByIdsRow
	ProviderMetadata interface{}
}
