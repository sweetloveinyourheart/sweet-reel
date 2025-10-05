package processing

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/samber/do"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/ffmpeg"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
)

const (
	BatchSize                 = 1024
	KafkaVideoProcessingGroup = "video-processing"
	KafkaVideoUploadedTopic   = "video-uploaded"

	S3VideoProcessedBucket = "video-processed"
)

type VideoProcessManager struct {
	ctx   context.Context
	queue chan lo.Tuple2[context.Context, *kafka.ConsumedMessage]

	storageClient s3.S3Storage
	ff            ffmpeg.FFmpegInterface
}

func NewVideoProcessManager(ctx context.Context) (*VideoProcessManager, error) {
	storageClient, err := do.Invoke[s3.S3Storage](nil)
	if err != nil {
		return nil, err
	}

	kafkaClient, err := do.Invoke[*kafka.Client](nil)
	if err != nil {
		return nil, err
	}

	vsp := &VideoProcessManager{
		ctx:           ctx,
		queue:         make(chan lo.Tuple2[context.Context, *kafka.ConsumedMessage], BatchSize*2),
		storageClient: storageClient,
		ff:            ffmpeg.New(),
	}

	go func() {
		messageHandler := func(ctx context.Context, msg *kafka.ConsumedMessage) error {
			logger.Global().InfoContext(ctx, "received message",
				zap.String("topic_name", msg.Topic), zap.String("key", msg.Key),
				zap.String("message",
					msg.ValueAsString()))

			switch msg.Topic {
			case KafkaVideoUploadedTopic:
				vsp.queue <- lo.T2(ctx, msg)
			}

			return nil
		}
		consumer, err := kafkaClient.CreateConsumer(KafkaVideoProcessingGroup, []string{KafkaVideoUploadedTopic}, messageHandler)
		if err != nil {
			logger.Global().ErrorContext(ctx, "failed to create consumer", zap.Error(err))
			return
		}

		err = consumer.Start(ctx)
		if err != nil {
			logger.Global().ErrorContext(ctx, "failed to start consumer", zap.Error(err))
			return
		}

		<-ctx.Done()
		consumer.Stop()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case tuple := <-vsp.queue:
				if err := vsp.HandleMessage(tuple.A, tuple.B); err != nil {
					logger.Global().ErrorContext(tuple.A, "failed to handle event", zap.Error(err))
				}
			}
		}
	}()

	return vsp, nil
}

func (vsp *VideoProcessManager) HandleMessage(ctx context.Context, message *kafka.ConsumedMessage) (err error) {
	if message == nil {
		return errors.Errorf("message is nil")
	}

	var msg s3.S3EventMessage
	if err := message.ValueAsJSON(&msg); err != nil {
		return err
	}

	bucket, key := s3.ExtractBucketAndKeyFromEventMessage(msg.Key)
	bytes, err := vsp.storageClient.Download(key, bucket)
	if err != nil {
		return err
	}

	// Process video using FFmpeg wrapper
	videoID := uuid.FromStringOrNil(s3.ExtractFileIDFromKey(msg.Key))
	if videoID == uuid.Nil {
		return errors.Errorf("invalid video id")
	}

	if err := vsp.processVideo(ctx, videoID, bytes); err != nil {
		return errors.Wrap(err, "failed to process video")
	}

	logger.Global().InfoContext(ctx, "Video processing completed successfully",
		zap.String("key", msg.Key))

	return nil
}

// processVideo handles the actual video processing using FFmpeg
func (vsp *VideoProcessManager) processVideo(ctx context.Context, videoID uuid.UUID, videoData []byte) error {
	if err := vsp.ff.IsAvailable(ctx); err != nil {
		return errors.Wrap(err, "FFmpeg not available")
	}

	tempDir := fmt.Sprintf("/tmp/video_processing_%s", videoID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create temp directory")
	}
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "input.mp4")
	if err := os.WriteFile(inputPath, videoData, 0644); err != nil {
		return errors.Wrap(err, "failed to write input file")
	}

	// Probe the input file to get information
	probeInfo, err := vsp.ff.ProbeFile(ctx, inputPath)
	if err != nil {
		return errors.Wrap(err, "failed to probe input file")
	}

	logger.Global().InfoContext(ctx, "Video file information",
		zap.String("format", probeInfo.Format.FormatName),
		zap.String("duration", probeInfo.Format.Duration),
		zap.String("size", probeInfo.Format.Size),
		zap.Int("streams", len(probeInfo.Streams)))

	hlsOutputDir := filepath.Join(tempDir, "hls")
	if err := os.MkdirAll(hlsOutputDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create HLS output directory")
	}

	// Define multiple quality levels for adaptive streaming
	qualities := []ffmpeg.SegmentationOptions{
		{
			SegmentDuration: "6",
			PlaylistType:    "vod",
			PlaylistName:    "playlist.m3u8",
			SegmentPrefix:   "segment",
			SegmentFormat:   "ts",
			VideoCodec:      "libx264",
			VideoBitrate:    "800k",
			AudioCodec:      "aac",
			AudioBitrate:    "96k",
			Resolution:      "854x480", // 480p
		},
		{
			SegmentDuration: "6",
			PlaylistType:    "vod",
			PlaylistName:    "playlist.m3u8",
			SegmentPrefix:   "segment",
			SegmentFormat:   "ts",
			VideoCodec:      "libx264",
			VideoBitrate:    "1400k",
			AudioCodec:      "aac",
			AudioBitrate:    "128k",
			Resolution:      "1280x720", // 720p
		},
		{
			SegmentDuration: "6",
			PlaylistType:    "vod",
			PlaylistName:    "playlist.m3u8",
			SegmentPrefix:   "segment",
			SegmentFormat:   "ts",
			VideoCodec:      "libx264",
			VideoBitrate:    "2800k",
			AudioCodec:      "aac",
			AudioBitrate:    "192k",
			Resolution:      "1920x1080", // 1080p
		},
	}

	progressCallback := func(progress ffmpeg.ProgressInfo) {
		if int(progress.Percentage)%10 == 0 { // Log every 10%
			logger.Global().InfoContext(ctx, "Video processing progress",
				zap.String("video_id", videoID.String()),
				zap.Float64("percentage", progress.Percentage),
				zap.String("speed", progress.Speed),
				zap.Duration("current", progress.Current),
				zap.Duration("total", progress.Duration))
		}
	}

	// Start video segmentation
	startTime := time.Now()
	logger.Global().InfoContext(ctx, "Starting video segmentation",
		zap.String("video_id", videoID.String()),
		zap.Int("quality_levels", len(qualities)))

	if err := vsp.ff.SegmentVideoMultiQuality(ctx, inputPath, hlsOutputDir, qualities, progressCallback); err != nil {
		return errors.Wrap(err, "failed to segment video")
	}

	processingTime := time.Since(startTime)
	logger.Global().InfoContext(ctx, "Video segmentation completed",
		zap.String("video_id", videoID.String()),
		zap.Duration("processing_time", processingTime))

	// Create thumbnail
	thumbnailPath := filepath.Join(tempDir, "thumbnail.jpg")
	if err := vsp.ff.CreateThumbnail(ctx, inputPath, thumbnailPath, "00:00:05", 320, 240); err != nil {
		logger.Global().WarnContext(ctx, "Failed to create thumbnail", zap.Error(err))
	} else {
		logger.Global().InfoContext(ctx, "Thumbnail created", zap.String("path", thumbnailPath))
	}

	// Upload processed files back to storage
	if err := vsp.uploadProcessedFiles(ctx, videoID, hlsOutputDir, thumbnailPath); err != nil {
		return errors.Wrap(err, "failed to upload processed files")
	}

	return nil
}

// uploadProcessedFiles uploads the HLS segments and thumbnail to storage
func (vsp *VideoProcessManager) uploadProcessedFiles(ctx context.Context, videoID uuid.UUID, hlsDir, thumbnailPath string) error {
	// Walk through HLS directory and upload all files
	err := filepath.Walk(hlsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Read file content
		fileData, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "failed to read file: %s", path)
		}

		// Generate storage key based on relative path
		relPath, err := filepath.Rel(hlsDir, path)
		if err != nil {
			return errors.Wrapf(err, "failed to get relative path for: %s", path)
		}

		storageKey := fmt.Sprintf("%s/hls/%s", videoID.String(), relPath)

		// Upload to storage
		fileReader := bytes.NewReader(fileData)
		mimeType := ffmpeg.GetMimeType(path)
		if err := vsp.storageClient.Upload(storageKey, S3VideoProcessedBucket, fileReader, mimeType); err != nil {
			return errors.Wrapf(err, "failed to upload file: %s", storageKey)
		}

		logger.Global().DebugContext(ctx, "Uploaded file to storage",
			zap.String("local_path", path),
			zap.String("storage_key", storageKey))

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to upload HLS files")
	}

	// Upload thumbnail if it exists
	if thumbnailPath != "" {
		if _, err := os.Stat(thumbnailPath); err == nil {
			thumbnailData, err := os.ReadFile(thumbnailPath)
			if err != nil {
				logger.Global().WarnContext(ctx, "Failed to read thumbnail", zap.Error(err))
			} else {
				thumbnailKey := fmt.Sprintf("%s/thumbnail.jpg", videoID.String())
				thumbnailReader := bytes.NewReader(thumbnailData)
				if err := vsp.storageClient.Upload(thumbnailKey, S3VideoProcessedBucket, thumbnailReader, "image/jpeg"); err != nil {
					logger.Global().WarnContext(ctx, "Failed to upload thumbnail", zap.Error(err))
				} else {
					logger.Global().InfoContext(ctx, "Thumbnail uploaded",
						zap.String("storage_key", thumbnailKey))
				}
			}
		}
	}

	logger.Global().InfoContext(ctx, "All processed files uploaded successfully",
		zap.String("video_id", videoID.String()))

	return nil
}
