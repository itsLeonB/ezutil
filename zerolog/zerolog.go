// Package zerolog provides Logger implementation that uses zerolog under the hood
package zerolog

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ungerr"
	"github.com/rs/zerolog"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

type zerologAdapter struct {
	logger zerolog.Logger
	// baseWriter is the original writer provided when creating the adapter.
	// It's used to create context-bound wrappers when WithContext() is called
	// so that per-adapter contexts are passed into writers that support it.
	baseWriter io.Writer
}

func NewZerologAdapter(writer io.Writer) *zerologAdapter {
	// Wrap the provided writer with a context-capable wrapper (initially nil
	// context). The wrapper will, when possible, call WriteWithContext on the
	// underlying writer so callers can pass an active context.
	wrapper := newCtxWrapper(writer, nil)

	return &zerologAdapter{
		logger:     zerolog.New(wrapper).With().Timestamp().Logger(),
		baseWriter: writer,
	}
}

// ctxWriter wraps an io.Writer and, if the wrapped writer implements
// WriteWithContext(ctx, p), will call that method with the provided context.
// Otherwise it will fall back to the plain Write method.
type ctxWriter struct {
	mu sync.RWMutex
	w  io.Writer
	// ctx can be nil
	ctx context.Context
}

func newCtxWrapper(w io.Writer, ctx context.Context) *ctxWriter {
	return &ctxWriter{w: w, ctx: ctx}
}

// SetContext sets the context used when writing. It's safe for concurrent use.
func (c *ctxWriter) SetContext(ctx context.Context) {
	c.mu.Lock()
	c.ctx = ctx
	c.mu.Unlock()
}

// Write implements io.Writer. If the wrapped writer implements a
// WriteWithContext(ctx, p) method, prefer calling that with the current
// context. Otherwise, delegate to the wrapped Write.
func (c *ctxWriter) Write(p []byte) (int, error) {
	c.mu.RLock()
	ctx := c.ctx
	w := c.w
	c.mu.RUnlock()

	type withCtx interface {
		WriteWithContext(context.Context, []byte) (int, error)
	}

	if wc, ok := w.(withCtx); ok {
		// if ctx is nil, still pass background to preserve previous behavior
		if ctx == nil {
			return wc.WriteWithContext(context.Background(), p)
		}
		return wc.WriteWithContext(ctx, p)
	}

	return w.Write(p)
}

// Debug logs a debug message.
func (z *zerologAdapter) Debug(args ...any) {
	z.logger.Debug().Msg(fmt.Sprint(args...))
}

// Info logs an info message.
func (z *zerologAdapter) Info(args ...any) {
	z.logger.Info().Msg(fmt.Sprint(args...))
}

// Warn logs a warning message.
func (z *zerologAdapter) Warn(args ...any) {
	z.logger.Warn().Msg(fmt.Sprint(args...))
}

// Error logs an error message.
func (z *zerologAdapter) Error(args ...any) {
	z.logger.Error().Msg(fmt.Sprint(args...))
}

// Fatal logs a fatal message and exits the program.
func (z *zerologAdapter) Fatal(args ...any) {
	z.logger.Fatal().Msg(fmt.Sprint(args...))
}

// Debugf logs a formatted debug message.
func (z *zerologAdapter) Debugf(format string, args ...any) {
	z.logger.Debug().Msg(fmt.Sprintf(format, args...))
}

// Infof logs a formatted info message.
func (z *zerologAdapter) Infof(format string, args ...any) {
	z.logger.Info().Msg(fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning message.
func (z *zerologAdapter) Warnf(format string, args ...any) {
	z.logger.Warn().Msg(fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error message.
func (z *zerologAdapter) Errorf(format string, args ...any) {
	z.logger.Error().Msg(fmt.Sprintf(format, args...))
}

// Fatalf logs a formatted fatal message and exits the program.
func (z *zerologAdapter) Fatalf(format string, args ...any) {
	z.logger.Fatal().Msg(fmt.Sprintf(format, args...))
}

// Print logs a message using Info level (goose.Logger interface).
func (z *zerologAdapter) Print(args ...any) { z.Info(args...) }

// Println logs a message using Info level (goose.Logger interface).
func (z *zerologAdapter) Println(args ...any) { z.Info(args...) }

// Printf logs a formatted message using Info level (goose.Logger interface).
func (z *zerologAdapter) Printf(format string, args ...any) { z.Infof(format, args...) }

func (z *zerologAdapter) WithError(err error) ezutil.Logger {
	if err == nil {
		return z
	}

	ctx := z.logger.With()

	switch e := err.(type) {
	case *ungerr.UnknownError:
		for _, attr := range e.ToLogAttrs() {
			ctx = ctx.Interface(attr.Key, attr.Value)
		}
	case ungerr.AppError:
		for _, attr := range e.ToLogAttrs() {
			ctx = ctx.Interface(attr.Key, attr.Value)
		}
	default:
		ctx = ctx.Str(string(semconv.ErrorTypeKey), fmt.Sprintf("%T", err))
		ctx = ctx.Str(string(semconv.ErrorMessageKey), err.Error())
	}

	return &zerologAdapter{logger: ctx.Logger(), baseWriter: z.baseWriter}
}

func (z *zerologAdapter) WithField(key string, value any) ezutil.Logger {
	return &zerologAdapter{logger: z.logger.With().Interface(key, value).Logger(), baseWriter: z.baseWriter}
}

func (z *zerologAdapter) WithFields(fields map[string]any) ezutil.Logger {
	ctx := z.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &zerologAdapter{logger: ctx.Logger(), baseWriter: z.baseWriter}
}

func (z *zerologAdapter) WithContext(ctx context.Context) ezutil.Logger {
	// Create a new wrapper that binds the provided context to the base writer
	// so writer implementations that accept a context can receive it.
	wrapper := newCtxWrapper(z.baseWriter, ctx)

	// Build a new zerolog logger that uses the context-bound writer.
	logger := zerolog.New(wrapper).With().Timestamp().Logger()

	lctx := logger.With()
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		sc := span.SpanContext()
		lctx = lctx.
			Str("trace_id", sc.TraceID().String()).
			Str("span_id", sc.SpanID().String())
	}

	return &zerologAdapter{logger: lctx.Logger(), baseWriter: z.baseWriter}
}

func Instance(l ezutil.Logger) zerolog.Logger {
	if adapter, ok := l.(*zerologAdapter); ok {
		return adapter.logger
	}
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
