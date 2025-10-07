package testing

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Model suite
type Model struct {
	suite.Suite
	*require.Assertions
}
