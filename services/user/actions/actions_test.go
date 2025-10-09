package actions_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	testingPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/repos"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/repos/mocks"
)

type ActionsSuite struct {
	*testingPkg.Suite
	mockUserRepository *mocks.MockUserRepository
}

func (as *ActionsSuite) SetupTest() {
	as.mockUserRepository = new(mocks.MockUserRepository)
}

func (as *ActionsSuite) TearDownTest() {
	as.mockUserRepository = nil
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

	do.Override(nil, func(i *do.Injector) (*pgxpool.Pool, error) {
		return nil, nil
	})
}
