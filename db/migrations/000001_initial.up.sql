CREATE TABLE
    IF NOT EXISTS sessions (
        id uuid DEFAULT gen_random_uuid (),
        file_size BIGINT NOT NULL,
        chunk_size BIGINT NOT NULL,
        max_chunk BIGINT NOT NULL,
        file_name TEXT,
        temp_path TEXT,
        owner_id TEXT,
        current_chunk BIGINT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS uploaded_files (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        file_name VARCHAR(255) NOT NULL,
        file_path TEXT NOT NULL,
        file_size BIGINT NOT NULL,
        file_type VARCHAR(100),
        uploaded_by TEXT NOT NULL,
        uploaded_at TIMESTAMP DEFAULT NOW (),
        is_deleted BOOLEAN DEFAULT FALSE,
        checksum VARCHAR(64),
        description TEXT,
        provider VARCHAR(50) NOT NULL,
        provider_metadata JSONB
    );

CREATE TABLE
    IF NOT EXISTS users (
        id uuid DEFAULT gen_random_uuid (),
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE 
    IF NOT EXISTS metrics (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    total_file_size BIGINT DEFAULT 0 NOT NULL,
    total_file_count BIGINT DEFAULT 0 NOT NULL,
    files_by_mimetype JSONB DEFAULT '{}'::JSONB,
    largest_file_size BIGINT DEFAULT 0 NOT NULL,
    smallest_file_size BIGINT DEFAULT 0 NOT NULL,
    average_file_size NUMERIC,
    total_deleted_files BIGINT DEFAULT 0 NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_uploaded_files_provider ON uploaded_files (provider);
CREATE INDEX idx_uploaded_files_size ON uploaded_files (file_size);
CREATE INDEX idx_uploaded_files_type ON uploaded_files (file_type);
CREATE INDEX idx_uploaded_files_duplicates ON uploaded_files (file_name, file_size, checksum);
CREATE INDEX idx_uploaded_files_provider_metadata ON uploaded_files USING gin (provider_metadata);

CREATE INDEX idx_owner_session ON sessions (owner_id);
