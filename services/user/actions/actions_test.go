package actions_test

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	testingPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/testing/mock"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/repos"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/repos/mocks"
)

type ActionsSuite struct {
	*testingPkg.Suite
	mockUserRepository    *mocks.MockUserRepository
	mockChannelRepository *mocks.MockChannelRepository
	mockConnPool          *mock.MockPgxPool
}

func (as *ActionsSuite) SetupTest() {
	as.mockUserRepository = new(mocks.MockUserRepository)
	as.mockChannelRepository = new(mocks.MockChannelRepository)
	as.mockConnPool = new(mock.MockPgxPool)
}

func (as *ActionsSuite) TearDownTest() {
	as.mockUserRepository = nil
	as.mockConnPool = nil
	as.mockChannelRepository = nil
}

func TestActionsSuite(t *testing.T) {
	as := &ActionsSuite{
		Suite: testingPkg.MakeSuite(t),
	}

	suite.Run(t, as)
}

func (as *ActionsSuite) setupEnvironment() {
	do.Override(nil, func(i *do.Injector) (repos.IUserRepository, error) {
		return as.mockUserRepository, nil
	})

	do.Override(nil, func(i *do.Injector) (repos.IChannelRepository, error) {
		return as.mockChannelRepository, nil
	})

	do.Override(nil, func(i *do.Injector) (db.ConnPool, error) {
		return as.mockConnPool, nil
	})
}
