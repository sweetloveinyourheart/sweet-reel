package processing

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/samber/do"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/jsonsutil"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/messages"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/models"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/repos"
)

const (
	BatchSize = 1024
)

type VideoProcessManager struct {
	ctx                 context.Context
	progressUpdateQueue chan lo.Tuple2[context.Context, *kafka.ConsumedMessage]
	videoProcessedQueue chan lo.Tuple2[context.Context, *kafka.ConsumedMessage]
	videoAggregateRepo  repos.IVideoAggregateRepository
}

func NewVideoProcessManager(ctx context.Context) (*VideoProcessManager, error) {
	kafkaClient, err := do.Invoke[*kafka.Client](nil)
	if err != nil {
		return nil, err
	}

	videoAggregateRepo, err := do.Invoke[repos.IVideoAggregateRepository](nil)
	if err != nil {
		logger.Global().Fatal("unable to get video aggregate repo")
	}

	vsp := &VideoProcessManager{
		ctx:                 ctx,
		progressUpdateQueue: make(chan lo.Tuple2[context.Context, *kafka.ConsumedMessage], BatchSize*2),
		videoProcessedQueue: make(chan lo.Tuple2[context.Context, *kafka.ConsumedMessage], BatchSize*2),
		videoAggregateRepo:  videoAggregateRepo,
	}

	go func() {
		messageHandler := func(ctx context.Context, msg *kafka.ConsumedMessage) error {
			logger.Global().InfoContext(ctx, "received message",
				zap.String("topic_name", msg.Topic), zap.String("key", msg.Key),
				zap.String("message",
					msg.ValueAsString()))

			switch msg.Topic {
			case kafka.KafkaVideoProgressTopic:
				vsp.progressUpdateQueue <- lo.T2(ctx, msg)
			case kafka.KafkaVideoProcessedTopic:
				vsp.videoProcessedQueue <- lo.T2(ctx, msg)
			}

			return nil
		}
		consumer, err := kafkaClient.CreateConsumer(kafka.KafkaVideoProcessingGroup,
			[]string{kafka.KafkaVideoProgressTopic, kafka.KafkaVideoProcessedTopic},
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
			case tuple := <-vsp.progressUpdateQueue:
				if err := vsp.HandleProgressUpdateMessage(tuple.A, tuple.B); err != nil {
					logger.Global().ErrorContext(tuple.A, "failed to handle video progress update event", zap.Error(err))
				}

			case tuple := <-vsp.videoProcessedQueue:
				if err := vsp.HandleVideoProcessedMessage(tuple.A, tuple.B); err != nil {
					logger.Global().ErrorContext(tuple.A, "failed to handle video processed event", zap.Error(err))
				}
			}
		}
	}()

	return vsp, nil
}

func (vsp *VideoProcessManager) HandleProgressUpdateMessage(ctx context.Context, message *kafka.ConsumedMessage) (err error) {
	if message == nil {
		return errors.New("message is nil")
	}

	var msg messages.VideoProcessingProgress
	if err := message.ValueAsJSON(&msg); err != nil {
		return err
	}

	err = vsp.videoAggregateRepo.UpdateVideoProgress(ctx,
		msg.VideoID,
		msg.ObjectKey,
		models.VideoStatus(msg.Status),
		msg.ProcessedAt,
	)
	if err != nil {
		return err
	}

	logger.Global().Info("video progress updated",
		zap.String("video_id", msg.VideoID.String()),
		zap.String("object_key", msg.ObjectKey),
		zap.String("status", string(msg.Status)))
	return nil
}

func (vsp *VideoProcessManager) HandleVideoProcessedMessage(ctx context.Context, message *kafka.ConsumedMessage) (err error) {
	if message == nil {
		return errors.New("message is nil")
	}

	var msg messages.VideoProcessed
	if err := message.ValueAsJSON(&msg); err != nil {
		return err
	}

	switch msg.Type {
	case messages.VideoProcessedTypeThumbnail:
		data, err := jsonsutil.ConvertData[messages.VideoProcessedThumbnailData](msg.Data)
		if err != nil {
			return errors.Wrap(err, "invalid thumbnail data")
		}

		newVideoThumbnail := &models.VideoThumbnail{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   msg.VideoID,
			ObjectKey: msg.ObjectKey,
			Width:     &data.Width,
			Height:    &data.Height,
		}
		err = vsp.videoAggregateRepo.CreateVideoThumbnail(ctx, newVideoThumbnail)
		if err != nil {
			return err
		}

	case messages.VideoProcessedTypeManifest:
		data, err := jsonsutil.ConvertData[messages.VideoProcessedManifestData](msg.Data)
		if err != nil {
			return errors.Wrap(err, "invalid manifest data")
		}

		newVideoManifest := &models.VideoManifest{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   msg.VideoID,
			ObjectKey: msg.ObjectKey,
			Quality:   data.Quality,
			SizeBytes: &data.SizeBytes,
		}
		err = vsp.videoAggregateRepo.CreateVideoManifest(ctx, newVideoManifest)
		if err != nil {
			return err
		}

	case messages.VideoProcessedTypeVariant:
		data, err := jsonsutil.ConvertData[messages.VideoProcessedVariantData](msg.Data)
		if err != nil {
			return errors.Wrap(err, "invalid video variant data")
		}

		newVideoManifest := &models.VideoVariant{
			ID:            uuid.Must(uuid.NewV7()),
			VideoID:       msg.VideoID,
			ObjectKey:     msg.ObjectKey,
			Quality:       data.Quality,
			TotalSegments: &data.TotalSegments,
			TotalDuration: &data.TotalDuration,
		}
		err = vsp.videoAggregateRepo.CreateVideoVariant(ctx, newVideoManifest)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid message type")
	}

	logger.Global().Info("processed files stored",
		zap.String("video_id", msg.VideoID.String()),
		zap.String("type", string(msg.Type)),
		zap.String("object_key", msg.ObjectKey))

	return nil
}
