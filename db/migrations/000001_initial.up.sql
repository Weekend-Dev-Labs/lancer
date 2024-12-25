CREATE TABLE
    IF NOT EXISTS sessions (
        id uuid DEFAULT gen_random_uuid (),
        file_size BIGINT NOT NULL,
        chunk_size BIGINT NOT NULL,
        max_chunk BIGINT NOT NULL,
        file_name TEXT,
        temp_path TEXT,
        owner_id TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_owner_session ON sessions (owner_id);