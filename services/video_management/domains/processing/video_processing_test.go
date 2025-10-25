package processing_test

import (
	"context"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	testingPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/repos"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/repos/mocks"
)

type VideoProcessingSuite struct {
	*testingPkg.Suite
	ctx    context.Context
	cancel context.CancelFunc

	mockVideoAggregateRepository *mocks.MockVideoAggregateRepository
}

func (as *VideoProcessingSuite) SetupTest() {
	as.mockVideoAggregateRepository = new(mocks.MockVideoAggregateRepository)
	as.ctx, as.cancel = context.WithTimeout(context.Background(), 10*time.Second)
}

func (as *VideoProcessingSuite) TearDownTest() {
	if as.cancel != nil {
		as.cancel()
	}

	as.mockVideoAggregateRepository = nil
}

func TestVideoProcessingSuite(t *testing.T) {
	as := &VideoProcessingSuite{
		Suite: testingPkg.MakeSuite(t),
	}

	suite.Run(t, as)
}

func (as *VideoProcessingSuite) setupEnvironment() {
	do.Override(nil, func(i *do.Injector) (*kafka.Client, error) {
		client := &kafka.Client{}
		return client, nil
	})

	do.Override(nil, func(i *do.Injector) (repos.IVideoAggregateRepository, error) {
		return as.mockVideoAggregateRepository, nil
	})
}
