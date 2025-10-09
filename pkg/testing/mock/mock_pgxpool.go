package mock

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

// MockPgxPool is a mock implementation that provides the same methods as *pgxpool.Pool
// Since we can't substitute for the concrete type directly, we'll modify the tests
// to work around this limitation
type MockPgxPool struct {
	mock.Mock
}

// AcquireFunc mocks the AcquireFunc method
func (m *MockPgxPool) AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

// Acquire mocks the Acquire method
func (m *MockPgxPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pgxpool.Conn), args.Error(1)
}

// QueryRow mocks the QueryRow method
func (m *MockPgxPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	mockArgs := m.Called(ctx, sql, args)
	if mockArgs.Get(0) == nil {
		return &MockRow{}
	}
	return mockArgs.Get(0).(pgx.Row)
}

// Query mocks the Query method
func (m *MockPgxPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mockArgs := m.Called(ctx, sql, args)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0).(pgx.Rows), mockArgs.Error(1)
}

// Exec mocks the Exec method
func (m *MockPgxPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	if args.Get(0) == nil {
		return pgconn.CommandTag{}, args.Error(1)
	}
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

// Begin mocks the Begin method
func (m *MockPgxPool) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}

// BeginTx mocks the BeginTx method
func (m *MockPgxPool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	args := m.Called(ctx, txOptions)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}

// Close mocks the Close method
func (m *MockPgxPool) Close() {
	m.Called()
}

// Ping mocks the Ping method
func (m *MockPgxPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Config mocks the Config method
func (m *MockPgxPool) Config() *pgxpool.Config {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgxpool.Config)
}

// Stat mocks the Stat method
func (m *MockPgxPool) Stat() *pgxpool.Stat {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgxpool.Stat)
}

// MockRow is a mock implementation of pgx.Row
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...any) error {
	args := m.Called(dest)
	return args.Error(0)
}

// Helper functions
func NewMockPgxPool() *MockPgxPool {
	return &MockPgxPool{}
}

func NewMockRow() *MockRow {
	return &MockRow{}
}
