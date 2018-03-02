package qingstor

import (
	"context"
	"io"

	"github.com/yunify/qscamel/model"
)

// Reachable implement source.Reachable
func (q *QingStor) Reachable() bool {
	return true
}

// Readable implement source.Readable
func (q *QingStor) Readable() bool {
	return true
}

// List implement source.List
func (q *QingStor) List(ctx context.Context, p string) (o []model.Object, err error) {
	return
}

// Read implement source.Read
func (q *QingStor) Read(ctx context.Context, p string) (r io.Reader, err error) {
	return
}
