package processing

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/samber/do"
	"github.com/samber/lo"
	"go.uber.org/zap"

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
	ctx       context.Context
	queue     chan lo.Tuple2[context.Context, *kafka.ConsumedMessage]
	videoRepo repos.IVideoRepository
}

func NewVideoProcessManager(ctx context.Context) (*VideoProcessManager, error) {
	kafkaClient, err := do.Invoke[*kafka.Client](nil)
	if err != nil {
		return nil, err
	}

	videoRepo, err := do.Invoke[repos.IVideoRepository](nil)
	if err != nil {
		logger.Global().Fatal("unable to get s3 client")
	}

	vsp := &VideoProcessManager{
		ctx:       ctx,
		queue:     make(chan lo.Tuple2[context.Context, *kafka.ConsumedMessage], BatchSize*2),
		videoRepo: videoRepo,
	}

	go func() {
		messageHandler := func(ctx context.Context, msg *kafka.ConsumedMessage) error {
			logger.Global().InfoContext(ctx, "received message",
				zap.String("topic_name", msg.Topic), zap.String("key", msg.Key),
				zap.String("message",
					msg.ValueAsString()))

			switch msg.Topic {
			case kafka.KafkaVideoProgressTopic:
				vsp.queue <- lo.T2(ctx, msg)
			}

			return nil
		}
		consumer, err := kafkaClient.CreateConsumer(kafka.KafkaVideoProcessingGroup,
			[]string{kafka.KafkaVideoProgressTopic},
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
		return errors.New("message is nil")
	}

	var msg messages.VideoProcessingProgress
	if err := message.ValueAsJSON(&msg); err != nil {
		return err
	}

	err = vsp.videoRepo.UpdateVideoProgress(ctx,
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
