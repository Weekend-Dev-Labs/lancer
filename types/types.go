package types

type ContextKeys string

const (
	ContextAuthInfo    = ContextKeys("auth-info")
	ContextSessionInfo = ContextKeys("session-info")
	ContextWebToken    = ContextKeys("webtoken-info")
)

type CreateSessionPayload struct {
	FileSize  int64  `json:"file_size" validate:"required"`
	ChunkSize int64  `json:"chunk_size" validate:"required"`
	MaxChunk  int64  `json:"max_chunk" validate:"required"`
	FileName  string `json:"file_name" validate:"required"`
}

type SessionTokenPayload struct {
	SessionToken string `json:"session_token" validate:"required"`
}

type SessionInfo struct {
	ID           string `json:"id"`
	FileSize     int64  `json:"file_size"`
	ChunkSize    int64  `json:"chunk_size"`
	MaxChunk     int64  `json:"max_chunk"`
	FileName     string `json:"file_name"`
	TempPath     string `json:"temp_path"`
	OwnerID      string `json:"owner_id"`
	CurrentChunk int64  `json:"current_chunk"`
}

type UploaderChunkPayload struct {
	ChunkCount int64  `validate:"required"`
	Checksum   string `validate:"required"`
}

type AdminUserPayload struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}
