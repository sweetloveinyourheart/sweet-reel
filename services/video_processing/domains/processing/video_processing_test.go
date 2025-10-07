package processing_test

import (
	"context"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	testingPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing"
	mockPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing/mock"
)

type VideoProcessingSuite struct {
	*testingPkg.Suite
	mockS3 *mockPkg.MockS3
	ctx    context.Context
	cancel context.CancelFunc
}

func (as *VideoProcessingSuite) SetupTest() {
	as.mockS3 = new(mockPkg.MockS3)
	as.ctx, as.cancel = context.WithTimeout(context.Background(), 10*time.Second)
}

func (as *VideoProcessingSuite) TearDownTest() {
	if as.cancel != nil {
		as.cancel()
	}
	as.mockS3 = nil
}

func TestVideoProcessingSuite(t *testing.T) {
	as := &VideoProcessingSuite{
		Suite: testingPkg.MakeSuite(t),
	}

	suite.Run(t, as)
}

func (as *VideoProcessingSuite) setupEnvironment() {
	do.Override(nil, func(i *do.Injector) (s3.S3Storage, error) {
		return as.mockS3, nil
	})

	do.Override(nil, func(i *do.Injector) (*kafka.Client, error) {
		client := &kafka.Client{}
		return client, nil
	})
}
