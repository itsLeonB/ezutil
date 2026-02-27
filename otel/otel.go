// Package otel implements instrumented zerolog logger
package otel

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/itsLeonB/ezutil/v2"
	ezerolog "github.com/itsLeonB/ezutil/v2/zerolog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
)

// internal/ezlog/otel.go
type otelWriter struct {
	logger log.Logger
}

func (w *otelWriter) Write(p []byte) (int, error) {
	// Keep backward-compatible Write behavior: delegate to WriteWithContext
	_, _ = w.WriteWithContext(context.Background(), p)
	return len(p), nil
}

// WriteWithContext emits the log record using the provided context so that
// trace/span correlation is preserved when available. This is the preferred
// method for callers that have an active context to propagate.
func (w *otelWriter) WriteWithContext(ctx context.Context, p []byte) (int, error) {
	record := log.Record{}
	record.SetTimestamp(time.Now())
	record.SetBody(log.StringValue(string(p)))
	w.logger.Emit(ctx, record)
	return len(p), nil
}

func Init(appNamespace string) ezutil.Logger {
	otelLogger := global.Logger(appNamespace)

	// Use a context-capable multi-writer so that when WithContext is used on
	// the zerolog adapter the active context flows into writers that support
	// WriteWithContext (such as otelWriter). This replaces zerolog.MultiLevelWriter
	// which doesn't expose WriteWithContext and would cause context to be lost.
	cw := &ctxMultiWriter{writers: []io.Writer{
		os.Stdout,
		&otelWriter{logger: otelLogger},
	}}

	return ezerolog.NewZerologAdapter(cw)
}

// ctxMultiWriter is like io.MultiWriter but also supports WriteWithContext so
// writers that accept a context can receive it when available.
type ctxMultiWriter struct {
	writers []io.Writer
}

func (m *ctxMultiWriter) Write(p []byte) (int, error) {
	var firstErr error
	for _, w := range m.writers {
		n, err := w.Write(p)
		if err == nil && n < len(p) && firstErr == nil {
			firstErr = io.ErrShortWrite
			continue
		}
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil {
		return 0, firstErr
	}
	return len(p), nil
}

func (m *ctxMultiWriter) WriteWithContext(ctx context.Context, p []byte) (int, error) {
	type withCtx interface {
		WriteWithContext(context.Context, []byte) (int, error)
	}

	var firstErr error
	for _, w := range m.writers {
		if wc, ok := w.(withCtx); ok {
			n, err := wc.WriteWithContext(ctx, p)
			if err == nil && n < len(p) && firstErr == nil {
				firstErr = io.ErrShortWrite
				continue
			}
			if err != nil && firstErr == nil {
				firstErr = err
			}
			continue
		}
		n, err := w.Write(p)
		if err == nil && n < len(p) && firstErr == nil {
			firstErr = io.ErrShortWrite
			continue
		}
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil {
		return 0, firstErr
	}
	return len(p), nil
}
