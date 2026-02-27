// Package otel implements instrumented zerolog logger
package otel

import (
	"context"
	"os"
	"time"

	"github.com/itsLeonB/ezutil/v2"
	ezerolog "github.com/itsLeonB/ezutil/v2/zerolog"
	"github.com/rs/zerolog"
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

	return ezerolog.NewZerologAdapter(zerolog.MultiLevelWriter(
		os.Stdout,
		&otelWriter{logger: otelLogger},
	))
}
