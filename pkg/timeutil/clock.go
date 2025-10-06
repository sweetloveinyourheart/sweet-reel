package timeutil

import "github.com/benbjohnson/clock"

var Clock clock.Clock
var MockedClock *clock.Mock

func init() {
	Clock = clock.New()
}

// MockClock sets the MockedClock field to a mock and sets the Clock field to the mock instance
func MockClock() {
	MockedClock = clock.NewMock()
	Clock = MockedClock
}
