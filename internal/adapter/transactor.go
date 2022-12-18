package adapter

import (
	"context"
	"fmt"
)

type (
	Transactor interface {
		Begin(ctx context.Context) error
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
)

var (
	ErrTxNotStarted = fmt.Errorf("tx not started")
)
