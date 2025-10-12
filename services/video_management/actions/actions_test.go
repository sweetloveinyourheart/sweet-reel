package actions_test

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	testingPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing"
	mockPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing/mock"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/repos"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/repos/mocks"
)

type ActionsSuite struct {
	*testingPkg.Suite
	mockVideoRepository *mocks.MockVideoRepository
	mockS3              *mockPkg.MockS3
}

func (as *ActionsSuite) SetupTest() {
	as.mockS3 = new(mockPkg.MockS3)
	as.mockVideoRepository = new(mocks.MockVideoRepository)
}

func (as *ActionsSuite) TearDownTest() {
	as.mockVideoRepository = nil
	as.mockS3 = nil
}

func TestActionsSuite(t *testing.T) {
	as := &ActionsSuite{
		Suite: testingPkg.MakeSuite(t),
	}

	suite.Run(t, as)
}

func (as *ActionsSuite) setupEnvironment() {
	do.Override(nil, func(i *do.Injector) (repos.IVideoRepository, error) {
		return as.mockVideoRepository, nil
	})

	do.Override(nil, func(i *do.Injector) (s3.S3Storage, error) {
		return as.mockS3, nil
	})
}
