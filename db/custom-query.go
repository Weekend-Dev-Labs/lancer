package db

import (
	"context"
	"fmt"
	"math"
)

type FilterUploadParams struct {
	MaxFileSize *int64
	MinFileSize *int64
	FileType    *string
	Checksum    *string
	UploadedBy  *string
	Provider    *string
	Limit       *int32
	Offset      *int32
}

type Metadata struct {
	Page       int32
	Size       int32
	TotalPage  int
	TotalCount int
}

type FilterUploadResult struct {
	Files []UploadedFile
	Meta  Metadata
}

func (db *Queries) GetFilteredUploads(params FilterUploadParams) (*FilterUploadResult, error) {

	data, err := db.FilterUploadedFiles(context.TODO(), FilterUploadedFilesParams{
		Column1: params.MinFileSize,
		Column2: params.MaxFileSize,
		Column3: params.FileType,
		Column4: params.Checksum,
		Column5: params.UploadedBy,
		Column6: params.Provider,
		Limit:   params.Limit,
		Offset:  params.Offset,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get data from database %v", err)
	}

	count, err := db.CountUploadedFiles(context.TODO(), CountUploadedFilesParams{
		Column1: params.MinFileSize,
		Column2: params.MaxFileSize,
		Column3: params.FileType,
		Column4: params.Checksum,
		Column5: params.UploadedBy,
		Column6: params.Provider,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to count data from database %v", err)
	}

	totalPages := float64(count / int64(*params.Limit))

	return &FilterUploadResult{
		Files: data,
		Meta: Metadata{
			Page:       *params.Offset,
			Size:       *params.Limit,
			TotalPage:  int(math.Ceil(float64(totalPages))),
			TotalCount: int(count),
		},
	}, nil
}

const filterUploadedFiles = `-- name: FilterUploadedFiles :many
SELECT 
    uf.id,
    uf.file_name,
    uf.file_path,
    uf.file_size,
    uf.file_type,
    uf.uploaded_by,
    uf.uploaded_at,
    uf.is_deleted,
    uf.checksum,
    uf.description,
    uf.provider,
    uf.provider_metadata
FROM uploaded_files uf
WHERE 
    uf.is_deleted = FALSE -- Always exclude deleted files
    AND ($1::BIGINT IS NULL OR uf.file_size >= $1)  -- Minimum file size
    AND ($2::BIGINT IS NULL OR uf.file_size <= $2)  -- Maximum file size
    AND ($3::TEXT IS NULL OR uf.file_type = $3)    -- File type filter
    AND ($4::TEXT IS NULL OR uf.checksum = $4)     -- Checksum filter
    AND ($5::TEXT IS NULL OR uf.uploaded_by = $5)  -- Uploaded by filter
    AND ($6::TEXT IS NULL OR uf.provider = $6)     -- Provider filter
ORDER BY 
    uf.uploaded_at DESC -- Sort by most recent uploads
LIMIT $7 -- Number of rows to return
OFFSET $8
`

type FilterUploadedFilesParams struct {
	Column1 *int64
	Column2 *int64
	Column3 *string
	Column4 *string
	Column5 *string
	Column6 *string
	Limit   *int32
	Offset  *int32
}

func (q *Queries) FilterUploadedFiles(ctx context.Context, arg FilterUploadedFilesParams) ([]UploadedFile, error) {
	rows, err := q.db.Query(ctx, filterUploadedFiles,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.Column4,
		arg.Column5,
		arg.Column6,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UploadedFile
	for rows.Next() {
		var i UploadedFile
		if err := rows.Scan(
			&i.ID,
			&i.FileName,
			&i.FilePath,
			&i.FileSize,
			&i.FileType,
			&i.UploadedBy,
			&i.UploadedAt,
			&i.IsDeleted,
			&i.Checksum,
			&i.Description,
			&i.Provider,
			&i.ProviderMetadata,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const countUploadedFiles = `-- name: CountUploadedFiles :one

SELECT 
    COUNT(*) AS total_count
FROM uploaded_files uf
WHERE 
    uf.is_deleted = FALSE -- Always exclude deleted files
    AND ($1::BIGINT IS NULL OR uf.file_size >= $1)  -- Minimum file size
    AND ($2::BIGINT IS NULL OR uf.file_size <= $2)  -- Maximum file size
    AND ($3::TEXT IS NULL OR uf.file_type = $3)    -- File type filter
    AND ($4::TEXT IS NULL OR uf.checksum = $4)     -- Checksum filter
    AND ($5::TEXT IS NULL OR uf.uploaded_by = $5)  -- Uploaded by filter
    AND ($6::TEXT IS NULL OR uf.provider = $6)
`

type CountUploadedFilesParams struct {
	Column1 *int64
	Column2 *int64
	Column3 *string
	Column4 *string
	Column5 *string
	Column6 *string
}

// Number of rows to skip
func (q *Queries) CountUploadedFiles(ctx context.Context, arg CountUploadedFilesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countUploadedFiles,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.Column4,
		arg.Column5,
		arg.Column6,
	)
	var total_count int64
	err := row.Scan(&total_count)
	return total_count, err
}
