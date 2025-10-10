package actions_test

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/oauth2"
	testingPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/testing/mock"
	userConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
)

type ActionsSuite struct {
	*testingPkg.Suite
	googleOAuthClient *mock.MockOAuthClient
	userServiceClient *mock.MockUserServiceClient
}

func (as *ActionsSuite) SetupTest() {
	as.googleOAuthClient = new(mock.MockOAuthClient)
	as.userServiceClient = new(mock.MockUserServiceClient)
}

func (as *ActionsSuite) TearDownTest() {
	as.googleOAuthClient = nil
	as.userServiceClient = nil
}

func TestActionsSuite(t *testing.T) {
	as := &ActionsSuite{
		Suite: testingPkg.MakeSuite(t),
	}

	suite.Run(t, as)
}

func (as *ActionsSuite) setupEnvironment() {
	do.OverrideNamed(nil, string(oauth2.ProviderGoogle), func(i *do.Injector) (oauth2.IOAuthClient, error) {
		return as.googleOAuthClient, nil
	})

	do.Override(nil, func(i *do.Injector) (userConnect.UserServiceClient, error) {
		return as.userServiceClient, nil
	})
}
