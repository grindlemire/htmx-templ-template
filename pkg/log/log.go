package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitGlobal initializes the log depending on the environment
func InitGlobal() {
	var core zapcore.Core

	var logger *zap.Logger
	if strings.ToLower(os.Getenv("ENV")) == "prod" {
		zapconf := zap.NewProductionConfig()
		zapconf.EncoderConfig.FunctionKey = "func"
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(zapconf.EncoderConfig),
			zapcore.Lock(os.Stdout),
			zapcore.InfoLevel,
		)

		logger = zap.New(
			core,
			zap.AddCaller(),
		).Named("prod")
	} else {
		// if we are not in gcp use a console logger
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
		}

		logLevel := zapcore.InfoLevel
		debug := strings.ToLower(os.Getenv("DEBUG"))
		if debug == "1" || debug == "true" {
			logLevel = zapcore.DebugLevel
		}

		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(config.EncoderConfig),
			zapcore.Lock(os.Stdout),
			logLevel,
		)

		logger = zap.New(
			core,
			zap.AddCaller(),
		).Named("dev")
	}

	zap.ReplaceGlobals(logger)
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type stackOpt struct {
	withStack bool
}

type opt func(*stackOpt)

func WithStack() opt {
	return func(o *stackOpt) {
		o.withStack = true
	}
}

func Callers(err error, opts ...opt) []zap.Field {
	o := &stackOpt{}
	for _, opt := range opts {
		opt(o)
	}

	stack := ""
	caller := ""
	fun := ""
	if err, ok := err.(stackTracer); ok {
		for i, f := range err.StackTrace() {
			if i == 0 {
				caller = fmt.Sprintf("%s:%d", f, f)
				fun = fmt.Sprintf("%n", f)
			}
			stack = fmt.Sprintf("%s%+s:%d\n", stack, f, f)
		}
	}

	fields := []zap.Field{
		zap.String("err_caller", caller),
		zap.String("err_func", fun),
	}

	if o.withStack {
		fields = append(fields, zap.String("err_stack", stack))
	}

	return fields
}
