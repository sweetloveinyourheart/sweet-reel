package repos

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/models"
)

type IVideoAggregateRepository interface {
	IVideoRepository
	GetVideosWithThumbnailByUploaderID(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]*models.VideoWithThumbnails, error)
}

type VideoAggregateRepository struct {
	VideoRepository
}

func NewVideoAggregateRepository(tx db.DbOrTx) IVideoAggregateRepository {
	return &VideoAggregateRepository{
		VideoRepository: VideoRepository{
			Tx: tx,
		},
	}
}

func (r *VideoRepository) GetVideosWithThumbnailByUploaderID(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]*models.VideoWithThumbnails, error) {
	query := `
		SELECT 
			videos.id, 
			uploader_id, 
			title, 
			description, 
			status, 
			videos.object_key, 
			processed_at, 
			videos.created_at, 
			videos.updated_at,
			video_thumbnails.id, 
			video_id, 
			video_thumbnails.object_key, 
			width,
			height, 
			video_thumbnails.created_at
		FROM videos 
		LEFT JOIN video_thumbnails ON videos.id = video_thumbnails.video_id
		WHERE uploader_id = $1 
		ORDER BY videos.created_at DESC, video_thumbnails.created_at ASC 
		LIMIT $2 OFFSET $3`

	rows, err := r.Tx.Query(ctx, query, uploaderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to group thumbnails by video ID
	videoMap := make(map[uuid.UUID]*models.VideoWithThumbnails)
	var videoOrder []uuid.UUID

	for rows.Next() {
		var (
			videoID            uuid.UUID
			uploaderID         uuid.UUID
			title              string
			description        *string
			status             models.VideoStatus
			objectKey          *string
			processedAt        *time.Time
			createdAt          time.Time
			updatedAt          time.Time
			thumbnailID        *uuid.UUID
			thumbnailVideoID   *uuid.UUID
			thumbnailObjectKey *string
			thumbnailWidth     *int
			thumbnailHeight    *int
			thumbnailCreatedAt *time.Time
		)

		err := rows.Scan(
			&videoID, &uploaderID, &title, &description,
			&status, &objectKey, &processedAt,
			&createdAt, &updatedAt,
			&thumbnailID, &thumbnailVideoID,
			&thumbnailObjectKey, &thumbnailWidth, &thumbnailHeight,
			&thumbnailCreatedAt)
		if err != nil {
			return nil, err
		}

		// Check if video already exists in map
		video, exists := videoMap[videoID]
		if !exists {
			// Create new video entry
			video = &models.VideoWithThumbnails{
				Video: models.Video{
					ID:          videoID,
					UploaderID:  uploaderID,
					Title:       title,
					Description: description,
					Status:      status,
					ObjectKey:   objectKey,
					ProcessedAt: processedAt,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
				Thumbnails: []models.VideoThumbnail{},
			}
			videoMap[videoID] = video
			videoOrder = append(videoOrder, videoID)
		}

		// Add thumbnail if it exists (LEFT JOIN may return NULL thumbnails)
		if thumbnailID != nil && thumbnailObjectKey != nil {
			thumbnail := models.VideoThumbnail{
				ID:        *thumbnailID,
				VideoID:   *thumbnailVideoID,
				ObjectKey: *thumbnailObjectKey,
				Width:     thumbnailWidth,
				Height:    thumbnailHeight,
				CreatedAt: *thumbnailCreatedAt,
			}
			video.Thumbnails = append(video.Thumbnails, thumbnail)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice preserving order
	videos := make([]*models.VideoWithThumbnails, 0, len(videoOrder))
	for _, videoID := range videoOrder {
		videos = append(videos, videoMap[videoID])
	}

	return videos, nil
}
