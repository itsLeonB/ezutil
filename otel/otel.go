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
	record := log.Record{}
	record.SetTimestamp(time.Now())
	record.SetBody(log.StringValue(string(p)))
	w.logger.Emit(context.Background(), record)
	return len(p), nil
}

func Init(appNamespace string) ezutil.Logger {
	otelLogger := global.Logger(appNamespace)

	return ezerolog.NewZerologAdapter(appNamespace, zerolog.MultiLevelWriter(
		os.Stdout,
		&otelWriter{logger: otelLogger},
	))
}
