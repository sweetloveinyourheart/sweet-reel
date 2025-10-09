package db

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/cockroachdb/errors"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.1"
	"go.opentelemetry.io/otel/trace"
)

var (
	dbTransactionScope = attribute.Key("db.transaction_ident")
	dbTransactionID    = attribute.Key("db.transaction_id")
	dbTransactionPID   = attribute.Key("db.transaction_pid")
)

type TransactionOption func(*TransactionOptions)

type TransactionOptions struct {
	ForceTransaction bool
	ForceSavepoint   bool
	IsolationLevel   sql.IsolationLevel
}

func ForceTransaction() TransactionOption {
	return func(options *TransactionOptions) {
		options.ForceTransaction = true
		options.ForceSavepoint = false
	}
}

func ForceSavepoint() TransactionOption {
	return func(options *TransactionOptions) {
		options.ForceTransaction = false
		options.ForceSavepoint = true
	}
}

func IsolationLevel(level sql.IsolationLevel) TransactionOption {
	return func(options *TransactionOptions) {
		options.IsolationLevel = level
	}
}

var defaultOptions = &TransactionOptions{}
var tracer trace.Tracer
var tracerOnce sync.Once

func init() {
	ResetDefaultOptions()
}

// ResetDefaultOptions resets default transaction options to the current configuration state
func ResetDefaultOptions() {
	defaultOptions.ForceSavepoint = false
	defaultOptions.ForceTransaction = true
}

func SetDefaultOptions(opts ...TransactionOption) {
	ResetDefaultOptions()
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(defaultOptions)
		}
	}

}

type TxInfo struct {
	TxID  any `db:"tx_id"`
	TxPID any `db:"tx_pid"`
}

// TransactionC starts a transaction or sets a save point on the connection.
// If the inner function returns an error then the transaction will be rolled back, otherwise the transaction will
// automatically commit at the end.
func TransactionC(ctx context.Context, db *pgx.Conn, scopeIdent string, fn func(tx pgx.Tx) error, opts ...TransactionOption) error {
	tracerOnce.Do(func() {
		tracer = otel.Tracer("com.sweetloveinyourheart.srl.pgx")
	})

	options := defaultOptions
	if len(opts) > 0 {
		options = &TransactionOptions{
			ForceTransaction: defaultOptions.ForceTransaction,
			ForceSavepoint:   defaultOptions.ForceSavepoint,
			IsolationLevel:   defaultOptions.IsolationLevel,
		}
		for _, opt := range opts {
			opt(options)
		}
	}
	if options.IsolationLevel != 0 {
		return PGXTransaction(db, SafeTransaction(ctx, scopeIdent, fn), &sql.TxOptions{
			Isolation: options.IsolationLevel,
		})
	}
	return PGXTransaction(db, SafeTransaction(ctx, scopeIdent, fn))

}

// Rollback starts a transaction or sets a save point on the connection.
// If the inner function returns an error then the transaction will be rolled back, otherwise the transaction will
// automatically rollback at the end.
func Rollback(ctx context.Context, db pgx.Tx, scopeIdent string, fn func(tx pgx.Tx)) error {
	return PGXTransaction(db, SafeTransaction(ctx, scopeIdent, func(tx pgx.Tx) error {
		fn(tx)
		return fmt.Errorf("auto rollback")
	}))
}

// SafeTransaction ensures that transactions that panic do so gracefully with rollback
func SafeTransaction(originalCtx context.Context, scopeIdent string, f func(tx pgx.Tx) error) func(tx pgx.Tx) error {
	return func(unboundTx pgx.Tx) error {
		txInfo, err := GetTXInfoPG(unboundTx)
		if err != nil {
			return err
		}

		attrs := []attribute.KeyValue{
			semconv.DBSystemPostgreSQL,
			dbTransactionScope.String(scopeIdent),
			dbTransactionID.String(fmt.Sprint(txInfo.TxID)),
			dbTransactionPID.String(fmt.Sprint(txInfo.TxPID)),
		}
		opts := []trace.SpanStartOption{
			trace.WithAttributes(attrs...),
			trace.WithSpanKind(trace.SpanKindClient),
		}
		spanCtx, span := tracer.Start(originalCtx, "Transaction", opts...)
		defer span.End()

		ctx := context.WithValue(spanCtx, dbTransactionID, txInfo.TxID)
		ctx = context.WithValue(ctx, dbTransactionPID, txInfo.TxPID)
		ctx = context.WithValue(ctx, dbTransactionScope, scopeIdent)
		defer TXRecover(ctx, unboundTx, false, "scope", "transaction_ident")

		logger.GlobalSugared().Infof("[transactions] opened transaction=%d pid(cid)=%d scope=%s", txInfo.TxID, txInfo.TxPID, scopeIdent)
		result := f(unboundTx)
		if result == nil {
			logger.GlobalSugared().Infof("[transactions] closed transaction=%d pid(cid)=%d scope=%s", txInfo.TxID, txInfo.TxPID, scopeIdent)
		} else {
			logger.GlobalSugared().Infof("[transactions] closed transaction=%d pid(cid)=%d scope=%s error=%s", txInfo.TxID, txInfo.TxPID, scopeIdent, result)
		}
		return result
	}
}

// GetTXInfoPG returns the transaction id and pid
func GetTXInfoPG(tx pgx.Tx) (*TxInfo, error) {
	type TxInfoPG struct {
		TxID  int64 `db:"tx_id"`
		TxPID int64 `db:"tx_pid"`
	}
	txInfo := &TxInfoPG{}
	if err := tx.QueryRow(context.Background(), "SELECT txid_current() AS tx_id, pg_backend_pid() AS tx_pid;").Scan(&txInfo.TxID, &txInfo.TxPID); err != nil {
		return nil, err
	}
	return &TxInfo{
		TxID:  txInfo.TxID,
		TxPID: txInfo.TxPID,
	}, nil
}

// TXRecover handles logging and rolling back transactions that panic
func TXRecover(ctx context.Context, tx pgx.Tx, isSavePoint bool, scopeType, scopeKey string) {
	r := recover()
	var err error
	if r != nil { //catch
		switch t := r.(type) {
		case error:
			err = t
		case string:
			err = errors.New(t)
		default:
			err = errors.New(fmt.Sprint(t))
		}

		// _internal_tx_id may not exist in the context so safely try to get its value if it does exist
		maybeTXID := ctx.Value("_internal_tx_id")
		txID := "unknown"
		if maybeTXID != nil {
			txID = fmt.Sprint(maybeTXID)
		}

		// _internal_tx_pid may not exist in the context so safely try to get its value if it does exist
		maybeTXPID := ctx.Value("_internal_tx_pid")
		txPID := "unknown"
		if maybeTXPID != nil {
			txPID = fmt.Sprint(maybeTXPID)
		}

		// _internal_tx_pid may not exist in the context so safely try to get its value if it does exist
		maybeScope := ctx.Value(scopeKey)
		scope := "unknown"
		if maybeScope != nil {
			scope = fmt.Sprint(maybeScope)
		}

		logger.GlobalSugared().Warnf("panic in handler, rolling back transaction=%s pid(cid)=%s %s=%s stack=%s", txID, txPID, scopeType, scope, debug.Stack())

		// The Savepoint func handles rollbacks in tests
		if !isSavePoint && tx != nil {
			if rollbackError := tx.Rollback(context.Background()); rollbackError != nil {
				err = fmt.Errorf("error rolling back transaction: original: %w, rollback: %#v", err, rollbackError)
				logger.GlobalSugared().Error(err)
			}
		}
	}

	// This should never panic
	// _internal_tx_complete may not exist in the context so safely try to get its value if it does exist
	maybeCompleteMonitorFunc := ctx.Value("_internal_tx_complete")
	if maybeCompleteMonitorFunc != nil {
		if completeMonitorFunc, ok := maybeCompleteMonitorFunc.(context.CancelFunc); ok {
			func() {
				defer func() {
					rec := recover()
					logger.GlobalSugared().Warnf("panic in monitor check %v %s", rec, scopeType)
				}()
				completeMonitorFunc()
			}()
		}
	}

	if err != nil {
		panic(err)
	}
}

func PGXTransaction(db Txi, fc func(tx pgx.Tx) error, opts ...*sql.TxOptions) (err error) {
	panicked := true

	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		// Make sure to rollback when panic, Block error or Commit error
		if panicked || err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()

	if err = fc(tx); err == nil {
		panicked = false
		return tx.Commit(context.Background())
	}

	panicked = false
	return
}

type Txi interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}
