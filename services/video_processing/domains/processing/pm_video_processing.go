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
	"github.com/sweetloveinyourheart/sweet-reel/pkg/messages"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
)

const (
	// Batch processing
	BatchSize = 1024

	// Thumbnail settings
	ThumbnailWidth      = 320
	ThumbnailHeight     = 240
	ThumbnailTimeOffset = "00:00:00"
	ThumbnailFileName   = "thumbnail.jpg"

	// Directory and file naming
	TempDirPattern = "/tmp/video_processing_%s"
	InputFileName  = "input.mp4"
	HLSDirName     = "hls"
)

// Re-export commonly used ffmpeg constants for convenience
const (
	// Quality levels
	QualityDefault = ffmpeg.QualityDefault
	Quality480p    = ffmpeg.Quality480p
	Quality720p    = ffmpeg.Quality720p
	Quality1080p   = ffmpeg.Quality1080p
	QualityUnknown = ffmpeg.QualityUnknown

	// MIME types
	ThumbnailMimeType = ffmpeg.MimeTypeJPEG

	// Playlist settings
	PlaylistFileName = ffmpeg.DefaultPlaylistName
	SegmentPrefix    = ffmpeg.DefaultSegmentPrefix
	SegmentDuration  = "6" // Custom duration, different from ffmpeg default
	PlaylistType     = ffmpeg.PlaylistTypeVOD
	SegmentFormat    = ffmpeg.DefaultSegmentFormat

	// Codecs
	CodecH264 = ffmpeg.CodecLibX264
	CodecAAC  = ffmpeg.CodecAAC

	// File extensions
	ExtM3U8 = ffmpeg.ExtM3U8
	ExtTS   = ffmpeg.ExtTS
)

// Video quality configurations
const (
	// Resolutions
	Resolution480p  = ffmpeg.Resolution480p
	Resolution720p  = ffmpeg.Resolution720p
	Resolution1080p = ffmpeg.Resolution1080p

	// Video bitrates
	Bitrate480p  = ffmpeg.VideoBitrate480p
	Bitrate720p  = ffmpeg.VideoBitrate720p
	Bitrate1080p = ffmpeg.VideoBitrate1080p

	// Audio bitrates
	AudioBitrate480p  = ffmpeg.AudioBitrate96k
	AudioBitrate720p  = ffmpeg.AudioBitrate128k
	AudioBitrate1080p = ffmpeg.AudioBitrate192k
)

var (
	// Define multiple quality levels for adaptive streaming
	qualities = []ffmpeg.SegmentationOptions{
		{
			QualityName:     Quality480p,
			SegmentDuration: SegmentDuration,
			PlaylistType:    PlaylistType,
			PlaylistName:    PlaylistFileName,
			SegmentPrefix:   SegmentPrefix,
			SegmentFormat:   SegmentFormat,
			VideoCodec:      CodecH264,
			VideoBitrate:    Bitrate480p,
			AudioCodec:      CodecAAC,
			AudioBitrate:    AudioBitrate480p,
			Resolution:      Resolution480p,
		},
		{
			QualityName:     Quality720p,
			SegmentDuration: SegmentDuration,
			PlaylistType:    PlaylistType,
			PlaylistName:    PlaylistFileName,
			SegmentPrefix:   SegmentPrefix,
			SegmentFormat:   SegmentFormat,
			VideoCodec:      CodecH264,
			VideoBitrate:    Bitrate720p,
			AudioCodec:      CodecAAC,
			AudioBitrate:    AudioBitrate720p,
			Resolution:      Resolution720p,
		},
		{
			QualityName:     Quality1080p,
			SegmentDuration: SegmentDuration,
			PlaylistType:    PlaylistType,
			PlaylistName:    PlaylistFileName,
			SegmentPrefix:   SegmentPrefix,
			SegmentFormat:   SegmentFormat,
			VideoCodec:      CodecH264,
			VideoBitrate:    Bitrate1080p,
			AudioCodec:      CodecAAC,
			AudioBitrate:    AudioBitrate1080p,
			Resolution:      Resolution1080p,
		},
	}
)

type VideoProcessManager struct {
	ctx   context.Context
	queue chan lo.Tuple2[context.Context, *kafka.ConsumedMessage]

	storageClient s3.S3Storage
	ff            ffmpeg.FFmpegInterface
	kafkaClient   *kafka.Client
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
		kafkaClient:   kafkaClient,
	}

	go func() {
		messageHandler := func(ctx context.Context, msg *kafka.ConsumedMessage) error {
			logger.Global().InfoContext(ctx, "received message",
				zap.String("topic_name", msg.Topic), zap.String("key", msg.Key),
				zap.String("message",
					msg.ValueAsString()))

			switch msg.Topic {
			case kafka.KafkaVideoUploadedTopic:
				vsp.queue <- lo.T2(ctx, msg)
			}

			return nil
		}
		consumer, err := kafkaClient.CreateConsumer(kafka.KafkaVideoProcessingGroup,
			[]string{kafka.KafkaVideoUploadedTopic},
			messageHandler)

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
		if err := consumer.Stop(); err != nil {
			logger.Global().ErrorContext(ctx, "failed to stop consumer", zap.Error(err))
		}
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

	var msg messages.S3EventMessage
	if err := message.ValueAsJSON(&msg); err != nil {
		return err
	}

	bucket, key := s3.ExtractBucketAndKeyFromEventMessage(msg.Key)
	bytes, err := vsp.storageClient.Download(key, bucket)
	if err != nil {
		return err
	}

	fileName, _ := s3.ExtractFilenameAndExt(key)
	videoID := uuid.FromStringOrNil(fileName)
	if videoID == uuid.Nil {
		return errors.Errorf("invalid video id: %s", videoID.String())
	}

	// Process video using FFmpeg wrapper
	err = vsp.processVideo(ctx, videoID, bytes)
	if err != nil {
		publishMsg := messages.VideoProcessingProgress{
			VideoID:     videoID,
			Status:      messages.VideoStatusFailed,
			ObjectKey:   key,
			ProcessedAt: time.Now(),
		}
		_, _, err := vsp.kafkaClient.SendJSON(ctx, kafka.KafkaVideoProgressTopic, videoID.String(), publishMsg)
		if err != nil {
			logger.Global().Error("Failed to publish video progress update message: %v", zap.Error(err))
			return err
		}

		return errors.Wrap(err, "failed to process video")
	} else {
		publishMsg := messages.VideoProcessingProgress{
			VideoID:     videoID,
			Status:      messages.VideoStatusReady,
			ObjectKey:   key,
			ProcessedAt: time.Now(),
		}
		_, _, err = vsp.kafkaClient.SendJSON(ctx, kafka.KafkaVideoProgressTopic, videoID.String(), publishMsg)
		if err != nil {
			logger.Global().Error("Failed to publish video progress update message: %v", zap.Error(err))
		}

		logger.Global().InfoContext(ctx, "Video processing completed successfully", zap.String("key", msg.Key))

		return nil
	}
}

// processVideo handles the actual video processing using FFmpeg
func (vsp *VideoProcessManager) processVideo(ctx context.Context, videoID uuid.UUID, videoData []byte) error {
	if err := vsp.ff.IsAvailable(ctx); err != nil {
		return errors.Wrap(err, "FFmpeg not available")
	}

	tempDir := fmt.Sprintf(TempDirPattern, videoID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create temp directory")
	}
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, InputFileName)
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

	hlsOutputDir := filepath.Join(tempDir, HLSDirName)
	if err := os.MkdirAll(hlsOutputDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create HLS output directory")
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
	thumbnailPath := filepath.Join(tempDir, ThumbnailFileName)
	if err := vsp.ff.CreateThumbnail(ctx, inputPath, thumbnailPath, ThumbnailTimeOffset, ThumbnailWidth, ThumbnailHeight); err != nil {
		logger.Global().WarnContext(ctx, "Failed to create thumbnail", zap.Error(err))
	} else {
		logger.Global().InfoContext(ctx, "Thumbnail created", zap.String("path", thumbnailPath))
	}

	// Upload processed files back to storage
	if err := vsp.uploadProcessedSegmentFiles(ctx, videoID, hlsOutputDir); err != nil {
		return errors.Wrap(err, "failed to upload processed segments files")
	}

	if err := vsp.uploadProcessedThumbnailFiles(ctx, videoID, thumbnailPath); err != nil {
		return errors.Wrap(err, "failed to upload processed thumbnail files")
	}

	return nil
}

// uploadProcessedFiles uploads the HLS segments to storage
func (vsp *VideoProcessManager) uploadProcessedSegmentFiles(ctx context.Context, videoID uuid.UUID, hlsDir string) error {
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
		if err := vsp.storageClient.Upload(storageKey, s3.S3VideoProcessedBucket, fileReader, mimeType); err != nil {
			return errors.Wrapf(err, "failed to upload file: %s", storageKey)
		}

		// Determine file type and publish appropriate message
		ext := filepath.Ext(path)
		switch ext {
		case ExtM3U8:
			quality := vsp.extractQualityFromPath(relPath)

			manifestData := messages.VideoProcessedManifestData{
				Quality:   quality,
				SizeBytes: info.Size(),
			}
			publishManifestMsg := messages.VideoProcessed{
				VideoID:   videoID,
				ObjectKey: storageKey,
				Type:      messages.VideoProcessedTypeManifest,
				Data:      manifestData,
			}
			_, _, err = vsp.kafkaClient.SendJSON(ctx, kafka.KafkaVideoProcessedTopic, videoID.String(), publishManifestMsg)
			if err != nil {
				logger.Global().Error("Failed to publish manifest message", zap.Error(err))
			}
		case ExtTS:
			// For variant segments, extract quality from path and count segments
			quality := vsp.extractQualityFromPath(relPath)
			segments := vsp.countSegmentsInDirectory(filepath.Dir(path))
			duration := vsp.calculateVariantDuration(filepath.Dir(path))

			variantData := messages.VideoProcessedVariantData{
				Quality:       quality,
				TotalSegments: segments,
				TotalDuration: duration,
			}
			publishVariantMsg := messages.VideoProcessed{
				VideoID:   videoID,
				ObjectKey: storageKey,
				Type:      messages.VideoProcessedTypeVariant,
				Data:      variantData,
			}
			_, _, err = vsp.kafkaClient.SendJSON(ctx, kafka.KafkaVideoProcessedTopic, videoID.String(), publishVariantMsg)
			if err != nil {
				logger.Global().Error("Failed to publish variant message", zap.Error(err))
			}
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to upload HLS files")
	}

	logger.Global().InfoContext(ctx, "All processed HLS files uploaded successfully",
		zap.String("video_id", videoID.String()))

	return nil
}

// uploadProcessedFiles uploads the thumbnail to storage
func (vsp *VideoProcessManager) uploadProcessedThumbnailFiles(ctx context.Context, videoID uuid.UUID, thumbnailPath string) error {
	// Upload thumbnail if it exists
	if thumbnailPath != "" {
		if _, err := os.Stat(thumbnailPath); err == nil {
			thumbnailData, err := os.ReadFile(thumbnailPath)
			if err != nil {
				logger.Global().WarnContext(ctx, "Failed to read thumbnail", zap.Error(err))
			} else {
				thumbnailKey := fmt.Sprintf("%s/%s", videoID.String(), ThumbnailFileName)
				thumbnailReader := bytes.NewReader(thumbnailData)
				if err := vsp.storageClient.Upload(thumbnailKey, s3.S3VideoProcessedBucket, thumbnailReader, ThumbnailMimeType); err != nil {
					logger.Global().WarnContext(ctx, "Failed to upload thumbnail", zap.Error(err))
				} else {
					data := messages.VideoProcessedThumbnailData{
						Width:  ThumbnailWidth,
						Height: ThumbnailHeight,
					}
					publishMsg := messages.VideoProcessed{
						VideoID:   videoID,
						ObjectKey: thumbnailKey,
						Type:      messages.VideoProcessedTypeThumbnail,
						Data:      data,
					}
					_, _, err := vsp.kafkaClient.SendJSON(ctx, kafka.KafkaVideoProcessedTopic, videoID.String(), publishMsg)
					if err != nil {
						logger.Global().Error("Failed to publish video progress update message: %v", zap.Error(err))
					}

					logger.Global().InfoContext(ctx, "Thumbnail uploaded", zap.String("storage_key", thumbnailKey))
				}
			}
		}
	}

	return nil
}

// extractQualityFromPath extracts the quality level from the file path
// e.g., "480p/segment_001.ts" -> "480p"
// e.g., "720p/segment_001.ts" -> "720p"
// e.g., "1080p/segment_001.ts" -> "1080p"
func (vsp *VideoProcessManager) extractQualityFromPath(relPath string) string {
	// Extract the quality directory from the path
	parts := filepath.SplitList(relPath)
	if len(parts) > 0 {
		dir := filepath.Dir(relPath)
		if dir == "." {
			return QualityDefault
		}
		return dir
	}

	return QualityUnknown
}

// countSegmentsInDirectory counts the number of .ts segment files in a directory
func (vsp *VideoProcessManager) countSegmentsInDirectory(dir string) int {
	count := 0
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ExtTS {
			count++
		}
	}

	return count
}

// calculateVariantDuration calculates the total duration of a variant by parsing its playlist
func (vsp *VideoProcessManager) calculateVariantDuration(dir string) int {
	playlistPath := filepath.Join(dir, PlaylistFileName)
	data, err := os.ReadFile(playlistPath)
	if err != nil {
		return 0
	}

	// Parse m3u8 file to calculate total duration
	// Look for #EXTINF tags which contain segment durations
	totalDuration := 0.0
	lines := bytes.SplitSeq(data, []byte("\n"))
	for line := range lines {
		lineStr := string(bytes.TrimSpace(line))
		if bytes.HasPrefix(line, []byte("#EXTINF:")) {
			// Extract duration from #EXTINF:6.000000,
			var segmentDuration float64
			_, err := fmt.Sscanf(lineStr, "#EXTINF:%f,", &segmentDuration)
			if err == nil {
				totalDuration += segmentDuration
			}
		}
	}

	return int(totalDuration)
}
