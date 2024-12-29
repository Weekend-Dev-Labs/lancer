CREATE TABLE
    uploaded_files (
        id SERIAL PRIMARY KEY,
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


CREATE INDEX idx_uploaded_files_provider ON uploaded_files (provider);

CREATE INDEX idx_uploaded_files_provider_metadata ON uploaded_files USING gin (provider_metadata);