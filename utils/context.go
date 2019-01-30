package utils

import (
	"context"
	"fmt"
)

// ContextKey is the type for context key.
type ContextKey string

// Context keys.
const (
	ContextKeyTask ContextKey = "task"
)

// FromTaskContext will extract task name from context.
func FromTaskContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(ContextKeyTask).(string); ok {
		fmt.Println(v)
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
