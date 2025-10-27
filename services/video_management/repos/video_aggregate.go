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
	GetChannelVideos(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]*models.ChannelVideo, error)
}

type VideoAggregateRepository struct {
	*VideoRepository
}

func NewVideoAggregateRepository(tx db.DbOrTx) IVideoAggregateRepository {
	return &VideoAggregateRepository{
		VideoRepository: &VideoRepository{
			Tx: tx,
		},
	}
}

func (r *VideoAggregateRepository) GetChannelVideos(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]*models.ChannelVideo, error) {
	query := `
		SELECT
			videos.id,
			uploader_id,
			channel_id,
			title,
			description,
			status,
			videos.object_key,
			processed_at,
			videos.created_at,
			videos.updated_at,
			videos.view_count,
			video_thumbnails.video_id,
			video_thumbnails.object_key,
			video_variants.video_id,
			video_variants.total_duration
		FROM videos
		LEFT JOIN video_thumbnails ON videos.id = video_thumbnails.video_id
		LEFT JOIN video_variants ON videos.id = video_variants.video_id
		WHERE channel_id = $1 AND status = 'ready'
		ORDER BY videos.created_at DESC, video_thumbnails.created_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.Tx.Query(ctx, query, channelID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to group thumbnails by video ID
	videoMap := make(map[uuid.UUID]*models.ChannelVideo)
	var videoOrder []uuid.UUID

	for rows.Next() {
		var (
			videoID              uuid.UUID
			uploaderID           uuid.UUID
			channelID            uuid.UUID
			title                string
			description          *string
			status               models.VideoStatus
			objectKey            *string
			processedAt          *time.Time
			createdAt            time.Time
			updatedAt            time.Time
			viewCount            int64
			thumbnailID          *uuid.UUID
			thumbnailObjectKey   *string
			variantID            *uuid.UUID
			variantTotalDuration *int
		)

		err := rows.Scan(
			&videoID, &uploaderID, &channelID, &title, &description,
			&status, &objectKey, &processedAt,
			&createdAt, &updatedAt, &viewCount,
			&thumbnailID, &thumbnailObjectKey,
			&variantID, &variantTotalDuration)
		if err != nil {
			return nil, err
		}

		// Check if video already exists in map
		video, exists := videoMap[videoID]
		if !exists {
			// Create new video entry
			video = &models.ChannelVideo{
				Video: models.Video{
					ID:          videoID,
					UploaderID:  uploaderID,
					ChannelID:   channelID,
					Title:       title,
					Description: description,
					Status:      status,
					ObjectKey:   objectKey,
					ProcessedAt: processedAt,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
				ThumbnailObjectKey: "",
				TotalDuration:      0,
				TotalView:          int(viewCount),
			}
			videoMap[videoID] = video
			videoOrder = append(videoOrder, videoID)
		}

		// Add thumbnail if it exists (LEFT JOIN may return NULL thumbnails)
		if thumbnailID != nil && thumbnailObjectKey != nil {
			video.ThumbnailObjectKey = *thumbnailObjectKey
		}

		// Add duration if it exists (LEFT JOIN may return NULL variants)
		if variantID != nil && variantTotalDuration != nil {
			video.TotalDuration = *variantTotalDuration
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice preserving order
	videos := make([]*models.ChannelVideo, 0, len(videoOrder))
	for _, videoID := range videoOrder {
		videos = append(videos, videoMap[videoID])
	}

	return videos, nil
}
