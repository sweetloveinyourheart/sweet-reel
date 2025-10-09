package db

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/jackc/pgx/v5"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/stringsutil"
)

var cnt uint64

// Savepoint will start a new save point in the transaction on the connection. If the inner function
// returns an error then the save point will be rolled back, otherwise the save point will not be rolled back to
// automatically at the end.
func Savepoint(tx pgx.Tx, scopeIdent string, fn func(tx pgx.Tx) error) error {
	//if committer, ok := tx.Statement.ConnPool.(gorm.TxCommitter); !ok || committer == nil {
	//	return errors.New("input connection is not inside a transaction")
	//}
	atomic.AddUint64(&cnt, 1)
	safeScope := fmt.Sprintf("subtx_%s_%d", stringsutil.ToMachineName(scopeIdent), cnt)
	defer func() {
		err := recover()
		if err != nil {
			_, dberr := tx.Exec(context.Background(), fmt.Sprintf("ROLLBACK TO %s;", safeScope))
			if dberr != nil {
				panic(fmt.Errorf("error committing or rolling back transaction. cause: %#v, txError: %w", err, dberr))
			}
			panic(err)
		}
	}()
	var dberr error

	_, err := tx.Exec(context.Background(), fmt.Sprintf("SAVEPOINT %s;", safeScope))
	if err != nil {
		return fmt.Errorf("error starting sub transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		_, dberr = tx.Exec(context.Background(), fmt.Sprintf("ROLLBACK TO %s;", safeScope))
		// PG Version
		//dberr = tx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT %s;", safeScope)).Error
	}

	if dberr != nil {
		return fmt.Errorf("error committing or rolling back transaction: %#v, source: %w", dberr, err)
	}

	return err
}
