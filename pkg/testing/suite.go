package testing

import (
	goTesting "testing"

	"github.com/stretchr/testify/require"
)

type Suite struct {
	*Model
}

func MakeSuite(t *goTesting.T) *Suite {
	model := &Model{}
	model.Assertions = require.New(t)

	result := &Suite{
		Model: model,
	}

	return result
}
