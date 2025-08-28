package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Env string

var (
	EnvProduction  Env = "production"
	EnvDevelopment Env = "development"
)

type Logger struct {
	zapLogger  *zap.Logger
	instanceId uuid.UUID
}

func New(
	environment Env,
	logLevel string,
	instanceId uuid.UUID,
) *Logger {
	if environment == EnvDevelopment {
		return &Logger{
			zapLogger:  newDevelopmentLogger(logLevel),
			instanceId: instanceId,
		}
	} else {
		return &Logger{
			zapLogger:  newProductionLogger(),
			instanceId: instanceId,
		}
	}
}

func ensureInit(l *Logger) *Logger {
	if l == nil {
		l = &Logger{
			zapLogger:  newProductionLogger(),
			instanceId: uuid.New(),
		}
	}
	return l
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	l = ensureInit(l)
	var fields []zap.Field

	fields = append(
		fields,
		zap.String(
			"instanceId",
			l.instanceId.String(),
		),
	)

	requestId, ok := ctx.Value("requestId").(string)
	if ok {
		fields = append(
			fields,
			zap.String(
				"requestId",
				requestId,
			),
		)
	}

	userId, ok := ctx.Value("userId").(int)
	if ok {
		fields = append(
			fields,
			zap.Int(
				"userId",
				userId,
			),
		)
	}

	l.zapLogger = l.zapLogger.With(fields...).Named("log")
	return l
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l = ensureInit(l)
	frames := getStackFrames()
	shortTrace := framesToShortString(frames)

	l.zapLogger.With(
		zap.Any(
			"caller",
			frames,
		),
		zap.String(
			"callerPath",
			shortTrace,
		),
	).Error(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l = ensureInit(l)
	l.zapLogger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l = ensureInit(l)
	l.zapLogger.Debug(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l = ensureInit(l)
	frames := getStackFrames()
	shortTrace := framesToShortString(frames)

	l.zapLogger.With(
		zap.Any(
			"caller",
			frames,
		),
		zap.String(
			"callerPath",
			shortTrace,
		),
	).Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l = ensureInit(l)
	frames := getStackFrames()
	shortTrace := framesToShortString(frames)

	l.zapLogger.With(
		zap.Any(
			"caller",
			frames,
		),
		zap.String(
			"callerPath",
			shortTrace,
		),
	).Fatal(msg, fields...)
}

func newDevelopmentLogger(logLvl string) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000"),
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Encoding:          "console",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		Level:             zap.NewAtomicLevelAt(mapLogLevel(logLvl)),
		Development:       true,
		DisableStacktrace: false,
	}

	loggerConfig, _ := zapConfig.Build()
	return loggerConfig
}

func newProductionLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000"),
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		zapcore.DebugLevel,
	)

	alertCore := sentryCore{LevelEnabler: zapcore.ErrorLevel}

	return zap.New(
		zapcore.NewTee(consoleCore, alertCore),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}

type sentryCore struct{ zapcore.LevelEnabler }

func (c sentryCore) With(fs []zapcore.Field) zapcore.Core { return c }

func (c sentryCore) Check(
	ent zapcore.Entry,
	ce *zapcore.CheckedEntry,
) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c sentryCore) Write(
	ent zapcore.Entry,
	fs []zapcore.Field,
) error {
	var err error
	var userId string
	var requestId string
	extras := map[string]any{}

	for _, f := range fs {
		switch f.Type {
		case zapcore.ErrorType:
			if e, ok := f.Interface.(error); ok {
				err = e
			}
		case zapcore.StringType:
			extras[f.Key] = f.String
			if f.Key == "clientId" {
				userId = f.String
			}
			if f.Key == "requestId" {
				requestId = f.String
			}
		case zapcore.BoolType:
			extras[f.Key] = f.Integer == 1
		case zapcore.Int64Type,
			zapcore.Int32Type,
			zapcore.Int8Type,
			zapcore.Int16Type,
			zapcore.Uint64Type,
			zapcore.Uint32Type,
			zapcore.Uint16Type,
			zapcore.Uint8Type:
			extras[f.Key] = f.Integer
			if f.Key == "clientId" {
				userId = strconv.FormatInt(f.Integer, 10)
			}
		default:
			extras[f.Key] = f.Interface
		}
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		if requestId != "" {
			scope.SetTag("requestId", requestId)
		}
		scope.SetTag("logSource", "zap")
		scope.SetExtras(extras)
		scope.SetLevel(sentry.LevelError)
		if userId != "" {
			scope.SetUser(sentry.User{
				ID: userId,
			})
		}
		captureWithStack(ent.Message, err)
	})

	return nil
}

func (c sentryCore) Sync() error {
	sentry.Flush(2 * time.Second)
	return nil
}

func mapLogLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func captureWithStack(msg string, err error) {
	const skip = 0

	pcs := pcsFromError(err)
	if len(pcs) == 0 {
		pcs = make([]uintptr, 64)
		n := runtime.Callers(skip, pcs)
		pcs = pcs[:n]
	}

	frames := runtime.CallersFrames(pcs)
	bi, ok := debug.ReadBuildInfo()
	packageName := "main"
	if ok {
		packageName = bi.Main.Path
	}

	var sentryFrames []sentry.Frame
	for {
		fr, more := frames.Next()
		if strings.Contains(
			fr.File,
			"logger.go",
		) || strings.Contains(
			fr.Function,
			"zapcore",
		) || strings.Contains(
			fr.Function, "sentry") {
			continue
		}
		sentryFr := sentry.NewFrame(fr)
		sentryFr.InApp = strings.HasPrefix(
			sentryFr.Module,
			packageName,
		)
		sentryFrames = append(
			[]sentry.Frame{sentryFr},
			sentryFrames...,
		)
		if !more {
			break
		}
	}

	ev := sentry.NewEvent()
	ev.Level = sentry.LevelError
	ev.Message = msg

	errType := msg
	errValue := ""
	if err != nil {
		errType = reflect.TypeOf(unwrapUntilNil(err)).String()
		errValue = err.Error()
	}

	ev.Exception = []sentry.Exception{{
		Type:       errType,
		Value:      errValue,
		Stacktrace: &sentry.Stacktrace{Frames: sentryFrames},
	}}

	sentry.CaptureEvent(ev)
}

func getStackFrames() []map[string]any {
	const maxDepth = 20
	pcs := make([]uintptr, maxDepth)

	n := runtime.Callers(0, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	var result []map[string]any
	for {
		frame, more := frames.Next()

		result = append([]map[string]any{{
			"package": frame.Function,
			"file":    filepath.Base(frame.File),
			"line":    frame.Line,
		}}, result...)

		if !more {
			break
		}
	}
	return result
}

func extractPackageName(funcFullName string) string {
	return funcFullName
}

func framesToShortString(frames []map[string]any) string {
	var parts []string
	for _, f := range frames {
		file := f["file"].(string)
		line := f["line"].(int)
		parts = append(parts, fmt.Sprintf("%s(%d)", file, line))
	}
	return strings.Join(parts, " -> ")
}

func pcsFromError(err error) []uintptr {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if stErr, ok := err.(stackTracer); ok {
		st := stErr.StackTrace()
		pcs := make([]uintptr, len(st))
		for i, fr := range st {
			pcs[i] = uintptr(fr)
		}
		return pcs
	}

	return nil
}

func unwrapUntilNil(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return errors.WithStack(err)
		}
		err = unwrapped
	}
}
