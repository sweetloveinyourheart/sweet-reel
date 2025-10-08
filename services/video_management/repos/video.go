package repos

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/models"
)

type VideoRepositoryInterface interface {
	// Video operations
	CreateVideo(ctx context.Context, video *models.Video) error
	GetVideoByID(ctx context.Context, id uuid.UUID) (*models.Video, error)
	GetVideosByUploaderID(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]*models.Video, error)
	UpdateVideo(ctx context.Context, video *models.Video) error
	UpdateVideoStatus(ctx context.Context, id uuid.UUID, status models.VideoStatus) error
	UpdateVideoProcessedAt(ctx context.Context, id uuid.UUID, processedAt time.Time) error
	DeleteVideo(ctx context.Context, id uuid.UUID) error
	ListVideos(ctx context.Context, limit, offset int) ([]*models.Video, error)

	// Video manifest operations
	CreateVideoManifest(ctx context.Context, manifest *models.VideoManifest) error
	GetVideoManifestByVideoID(ctx context.Context, videoID uuid.UUID) (*models.VideoManifest, error)
	UpdateVideoManifest(ctx context.Context, manifest *models.VideoManifest) error
	DeleteVideoManifest(ctx context.Context, id uuid.UUID) error

	// Video variant operations
	CreateVideoVariant(ctx context.Context, variant *models.VideoVariant) error
	GetVideoVariantsByVideoID(ctx context.Context, videoID uuid.UUID) ([]*models.VideoVariant, error)
	GetVideoVariantByID(ctx context.Context, id uuid.UUID) (*models.VideoVariant, error)
	UpdateVideoVariant(ctx context.Context, variant *models.VideoVariant) error
	DeleteVideoVariant(ctx context.Context, id uuid.UUID) error
	DeleteVideoVariantsByVideoID(ctx context.Context, videoID uuid.UUID) error

	// Video thumbnail operations
	CreateVideoThumbnail(ctx context.Context, thumbnail *models.VideoThumbnail) error
	GetVideoThumbnailsByVideoID(ctx context.Context, videoID uuid.UUID) ([]*models.VideoThumbnail, error)
	GetVideoThumbnailByID(ctx context.Context, id uuid.UUID) (*models.VideoThumbnail, error)
	UpdateVideoThumbnail(ctx context.Context, thumbnail *models.VideoThumbnail) error
	DeleteVideoThumbnail(ctx context.Context, id uuid.UUID) error
	DeleteVideoThumbnailsByVideoID(ctx context.Context, videoID uuid.UUID) error

	// Aggregate operations
	GetVideoWithRelations(ctx context.Context, videoID uuid.UUID) (*models.Video, *models.VideoManifest, []*models.VideoVariant, []*models.VideoThumbnail, error)
	GetVideoCount(ctx context.Context) (int64, error)
	GetVideoCountByUploaderID(ctx context.Context, uploaderID uuid.UUID) (int64, error)
}

type VideoRepository struct {
	Tx db.DbOrTx
}

func NewVideoRepository(tx db.DbOrTx) VideoRepositoryInterface {
	return &VideoRepository{
		Tx: tx,
	}
}

// Video operations

func (r *VideoRepository) CreateVideo(ctx context.Context, video *models.Video) error {
	query := `
		INSERT INTO videos (id, uploader_id, title, description, status, original_file_url, processed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.Tx.Exec(ctx, query,
		video.ID, video.UploaderID, video.Title, video.Description,
		video.Status, video.OriginalFileURL, video.ProcessedAt,
		video.CreatedAt, video.UpdatedAt)
	return err
}

func (r *VideoRepository) GetVideoByID(ctx context.Context, id uuid.UUID) (*models.Video, error) {
	query := `
		SELECT id, uploader_id, title, description, status, original_file_url, processed_at, created_at, updated_at
		FROM videos WHERE id = $1`

	video := &models.Video{}
	err := r.Tx.QueryRow(ctx, query, id).Scan(
		&video.ID, &video.UploaderID, &video.Title, &video.Description,
		&video.Status, &video.OriginalFileURL, &video.ProcessedAt,
		&video.CreatedAt, &video.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return video, nil
}

func (r *VideoRepository) GetVideosByUploaderID(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]*models.Video, error) {
	query := `
		SELECT id, uploader_id, title, description, status, original_file_url, processed_at, created_at, updated_at
		FROM videos WHERE uploader_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.Tx.Query(ctx, query, uploaderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*models.Video
	for rows.Next() {
		video := &models.Video{}
		err := rows.Scan(
			&video.ID, &video.UploaderID, &video.Title, &video.Description,
			&video.Status, &video.OriginalFileURL, &video.ProcessedAt,
			&video.CreatedAt, &video.UpdatedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, rows.Err()
}

func (r *VideoRepository) UpdateVideo(ctx context.Context, video *models.Video) error {
	query := `
		UPDATE videos SET uploader_id = $2, title = $3, description = $4, status = $5, 
		original_file_url = $6, processed_at = $7, updated_at = $8
		WHERE id = $1`

	_, err := r.Tx.Exec(ctx, query,
		video.ID, video.UploaderID, video.Title, video.Description,
		video.Status, video.OriginalFileURL, video.ProcessedAt, video.UpdatedAt)
	return err
}

func (r *VideoRepository) UpdateVideoStatus(ctx context.Context, id uuid.UUID, status models.VideoStatus) error {
	query := `UPDATE videos SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, id, status)
	return err
}

func (r *VideoRepository) UpdateVideoProcessedAt(ctx context.Context, id uuid.UUID, processedAt time.Time) error {
	query := `UPDATE videos SET processed_at = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, id, processedAt)
	return err
}

func (r *VideoRepository) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM videos WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, id)
	return err
}

func (r *VideoRepository) ListVideos(ctx context.Context, limit, offset int) ([]*models.Video, error) {
	query := `
		SELECT id, uploader_id, title, description, status, original_file_url, processed_at, created_at, updated_at
		FROM videos ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.Tx.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*models.Video
	for rows.Next() {
		video := &models.Video{}
		err := rows.Scan(
			&video.ID, &video.UploaderID, &video.Title, &video.Description,
			&video.Status, &video.OriginalFileURL, &video.ProcessedAt,
			&video.CreatedAt, &video.UpdatedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, rows.Err()
}

// Video manifest operations

func (r *VideoRepository) CreateVideoManifest(ctx context.Context, manifest *models.VideoManifest) error {
	query := `
		INSERT INTO video_manifests (id, video_id, manifest_url, size_bytes, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.Tx.Exec(ctx, query,
		manifest.ID, manifest.VideoID, manifest.ManifestURL,
		manifest.SizeBytes, manifest.CreatedAt)
	return err
}

func (r *VideoRepository) GetVideoManifestByVideoID(ctx context.Context, videoID uuid.UUID) (*models.VideoManifest, error) {
	query := `
		SELECT id, video_id, manifest_url, size_bytes, created_at
		FROM video_manifests WHERE video_id = $1`

	manifest := &models.VideoManifest{}
	err := r.Tx.QueryRow(ctx, query, videoID).Scan(
		&manifest.ID, &manifest.VideoID, &manifest.ManifestURL,
		&manifest.SizeBytes, &manifest.CreatedAt)

	if err != nil {
		return nil, err
	}
	return manifest, nil
}

func (r *VideoRepository) UpdateVideoManifest(ctx context.Context, manifest *models.VideoManifest) error {
	query := `
		UPDATE video_manifests SET video_id = $2, manifest_url = $3, size_bytes = $4
		WHERE id = $1`

	_, err := r.Tx.Exec(ctx, query,
		manifest.ID, manifest.VideoID, manifest.ManifestURL, manifest.SizeBytes)
	return err
}

func (r *VideoRepository) DeleteVideoManifest(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM video_manifests WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, id)
	return err
}

// Video variant operations

func (r *VideoRepository) CreateVideoVariant(ctx context.Context, variant *models.VideoVariant) error {
	query := `
		INSERT INTO video_variants (id, video_id, quality, playlist_url, total_segments, total_duration, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.Tx.Exec(ctx, query,
		variant.ID, variant.VideoID, variant.Quality, variant.PlaylistURL,
		variant.TotalSegments, variant.TotalDuration, variant.CreatedAt)
	return err
}

func (r *VideoRepository) GetVideoVariantsByVideoID(ctx context.Context, videoID uuid.UUID) ([]*models.VideoVariant, error) {
	query := `
		SELECT id, video_id, quality, playlist_url, total_segments, total_duration, created_at
		FROM video_variants WHERE video_id = $1 ORDER BY quality`

	rows, err := r.Tx.Query(ctx, query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []*models.VideoVariant
	for rows.Next() {
		variant := &models.VideoVariant{}
		err := rows.Scan(
			&variant.ID, &variant.VideoID, &variant.Quality, &variant.PlaylistURL,
			&variant.TotalSegments, &variant.TotalDuration, &variant.CreatedAt)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}
	return variants, rows.Err()
}

func (r *VideoRepository) GetVideoVariantByID(ctx context.Context, id uuid.UUID) (*models.VideoVariant, error) {
	query := `
		SELECT id, video_id, quality, playlist_url, total_segments, total_duration, created_at
		FROM video_variants WHERE id = $1`

	variant := &models.VideoVariant{}
	err := r.Tx.QueryRow(ctx, query, id).Scan(
		&variant.ID, &variant.VideoID, &variant.Quality, &variant.PlaylistURL,
		&variant.TotalSegments, &variant.TotalDuration, &variant.CreatedAt)

	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (r *VideoRepository) UpdateVideoVariant(ctx context.Context, variant *models.VideoVariant) error {
	query := `
		UPDATE video_variants SET video_id = $2, quality = $3, playlist_url = $4, 
		total_segments = $5, total_duration = $6 WHERE id = $1`

	_, err := r.Tx.Exec(ctx, query,
		variant.ID, variant.VideoID, variant.Quality, variant.PlaylistURL,
		variant.TotalSegments, variant.TotalDuration)
	return err
}

func (r *VideoRepository) DeleteVideoVariant(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM video_variants WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, id)
	return err
}

func (r *VideoRepository) DeleteVideoVariantsByVideoID(ctx context.Context, videoID uuid.UUID) error {
	query := `DELETE FROM video_variants WHERE video_id = $1`
	_, err := r.Tx.Exec(ctx, query, videoID)
	return err
}

// Video thumbnail operations

func (r *VideoRepository) CreateVideoThumbnail(ctx context.Context, thumbnail *models.VideoThumbnail) error {
	query := `
		INSERT INTO video_thumbnails (id, video_id, file_url, width, height, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.Tx.Exec(ctx, query,
		thumbnail.ID, thumbnail.VideoID, thumbnail.FileURL,
		thumbnail.Width, thumbnail.Height, thumbnail.CreatedAt)
	return err
}

func (r *VideoRepository) GetVideoThumbnailsByVideoID(ctx context.Context, videoID uuid.UUID) ([]*models.VideoThumbnail, error) {
	query := `
		SELECT id, video_id, file_url, width, height, created_at
		FROM video_thumbnails WHERE video_id = $1 ORDER BY created_at`

	rows, err := r.Tx.Query(ctx, query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var thumbnails []*models.VideoThumbnail
	for rows.Next() {
		thumbnail := &models.VideoThumbnail{}
		err := rows.Scan(
			&thumbnail.ID, &thumbnail.VideoID, &thumbnail.FileURL,
			&thumbnail.Width, &thumbnail.Height, &thumbnail.CreatedAt)
		if err != nil {
			return nil, err
		}
		thumbnails = append(thumbnails, thumbnail)
	}
	return thumbnails, rows.Err()
}

func (r *VideoRepository) GetVideoThumbnailByID(ctx context.Context, id uuid.UUID) (*models.VideoThumbnail, error) {
	query := `
		SELECT id, video_id, file_url, width, height, created_at
		FROM video_thumbnails WHERE id = $1`

	thumbnail := &models.VideoThumbnail{}
	err := r.Tx.QueryRow(ctx, query, id).Scan(
		&thumbnail.ID, &thumbnail.VideoID, &thumbnail.FileURL,
		&thumbnail.Width, &thumbnail.Height, &thumbnail.CreatedAt)

	if err != nil {
		return nil, err
	}
	return thumbnail, nil
}

func (r *VideoRepository) UpdateVideoThumbnail(ctx context.Context, thumbnail *models.VideoThumbnail) error {
	query := `
		UPDATE video_thumbnails SET video_id = $2, file_url = $3, width = $4, height = $5
		WHERE id = $1`

	_, err := r.Tx.Exec(ctx, query,
		thumbnail.ID, thumbnail.VideoID, thumbnail.FileURL,
		thumbnail.Width, thumbnail.Height)
	return err
}

func (r *VideoRepository) DeleteVideoThumbnail(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM video_thumbnails WHERE id = $1`
	_, err := r.Tx.Exec(ctx, query, id)
	return err
}

func (r *VideoRepository) DeleteVideoThumbnailsByVideoID(ctx context.Context, videoID uuid.UUID) error {
	query := `DELETE FROM video_thumbnails WHERE video_id = $1`
	_, err := r.Tx.Exec(ctx, query, videoID)
	return err
}

// Aggregate operations

func (r *VideoRepository) GetVideoWithRelations(ctx context.Context, videoID uuid.UUID) (*models.Video, *models.VideoManifest, []*models.VideoVariant, []*models.VideoThumbnail, error) {
	video, err := r.GetVideoByID(ctx, videoID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	manifest, _ := r.GetVideoManifestByVideoID(ctx, videoID)     // May not exist
	variants, _ := r.GetVideoVariantsByVideoID(ctx, videoID)     // May be empty
	thumbnails, _ := r.GetVideoThumbnailsByVideoID(ctx, videoID) // May be empty

	return video, manifest, variants, thumbnails, nil
}

func (r *VideoRepository) GetVideoCount(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM videos`
	var count int64
	err := r.Tx.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *VideoRepository) GetVideoCountByUploaderID(ctx context.Context, uploaderID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM videos WHERE uploader_id = $1`
	var count int64
	err := r.Tx.QueryRow(ctx, query, uploaderID).Scan(&count)
	return count, err
}
