package utils

import (
	"context"

	"github.com/yunify/qscamel/db"
)

// ContextKey is the type for context key.
type ContextKey string

// Context keys.
const (
	ContextKeyTx   ContextKey = "tx"
	ContextKeyTask ContextKey = "task"
)

// FromTxContext will extract tx from context.
func FromTxContext(ctx context.Context) *db.Tx {
	if ctx == nil {
		return nil
	}
	if v, ok := ctx.Value(ContextKeyTx).(*db.Tx); ok {
		return v
	}
	return nil
}

// NewTxContext will create a ctx with tx.
func NewTxContext(ctx context.Context, tx *db.Tx) context.Context {
	if ctx == nil || tx == nil {
		return ctx
	}
	// If ctx already has a tx, we will return ctx directly.
	if _, ok := ctx.Value(ContextKeyTx).(*db.Tx); ok {
		return ctx
	}
	return context.WithValue(ctx, ContextKeyTx, tx)
}

// FromTaskContext will extract task name from context.
func FromTaskContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(ContextKeyTask).(string); ok {
		return v
	}
	return ""
}

// NewTaskContext will create a ctx with task name.
func NewTaskContext(ctx context.Context, t string) context.Context {
	if ctx == nil || t == "" {
		return ctx
	}
	// If ctx already has a tx, we will return ctx directly.
	if _, ok := ctx.Value(ContextKeyTask).(string); ok {
		return ctx
	}
	return context.WithValue(ctx, ContextKeyTask, t)
}
