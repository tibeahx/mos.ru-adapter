package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const Source = "source"

var sugar = newSugar()

func newSugar() *zap.SugaredLogger {
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:          "T",
			LevelKey:         "L",
			NameKey:          "N",
			CallerKey:        "C",
			FunctionKey:      zapcore.OmitKey,
			MessageKey:       "M",
			StacktraceKey:    "S",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeDuration:   zapcore.StringDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return log.WithOptions(zap.AddStacktrace(zap.ErrorLevel)).Sugar()
}

func Zap() *zap.SugaredLogger {
	return sugar
}

type ctxKey struct{}

func CtxWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctxLogger, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return ctxLogger
	}

	return sugar.Desugar()
}

func WithSource(logger *zap.Logger, source string) *zap.Logger {
	return logger.With(zap.String(Source, source))
}
