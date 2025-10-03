package processing

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/samber/do"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/storage"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_processing/models"
)

const (
	BatchSize                 = 1024
	KafkaVideoProcessingGroup = "video_processing"
	KafkaVideoSplitterTopic   = "video_splitter"
)

type VideoSplitterProcessManager struct {
	ctx   context.Context
	queue chan lo.Tuple2[context.Context, *kafka.ConsumedMessage]

	storageClient storage.Storage
}

func NewVideoSplitterProcessManager(ctx context.Context) (*VideoSplitterProcessManager, error) {
	storageClient, err := do.Invoke[storage.Storage](nil)
	if err != nil {
		return nil, err
	}

	kafkaClient, err := do.Invoke[*kafka.Client](nil)
	if err != nil {
		return nil, err
	}

	vsp := &VideoSplitterProcessManager{
		ctx:           ctx,
		queue:         make(chan lo.Tuple2[context.Context, *kafka.ConsumedMessage], BatchSize*2),
		storageClient: storageClient,
	}

	go func() {
		messageHandler := func(ctx context.Context, msg *kafka.ConsumedMessage) error {
			logger.Global().InfoContext(ctx, "received message", zap.String("topic_name", msg.Topic), zap.String("key", msg.Key), zap.String("message", msg.ValueAsString()))

			switch msg.Topic {
			case KafkaVideoSplitterTopic:
				vsp.queue <- lo.T2(ctx, msg)
			}

			return nil
		}
		consumer, err := kafkaClient.CreateConsumer(KafkaVideoProcessingGroup, []string{KafkaVideoSplitterTopic}, messageHandler)
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

func (vsp *VideoSplitterProcessManager) HandleMessage(ctx context.Context, message *kafka.ConsumedMessage) (err error) {
	if message == nil {
		return errors.Errorf("message is nil")
	}

	var msg models.VideoSplitterMessage
	if err := message.ValueAsJSON(&msg); err != nil {
		return err
	}

	bytes, err := vsp.storageClient.Download(msg.Metadata.Key, msg.Metadata.Bucket)
	if err != nil {
		return err
	}

	// TODO: Handle splitter logic
	logger.Global().Info("Download successfully")
	_ = bytes

	return nil
}
