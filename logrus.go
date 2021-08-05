package logrus

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/unistack-org/micro/v3/logger"
)

type Logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type logrusLogger struct {
	Logger Logger
	opts   Options
}

func (l *logrusLogger) Init(opts ...logger.Option) error {
	for _, o := range opts {
		o(&l.opts.Options)
	}

	if formatter, ok := l.opts.Context.Value(formatterKey{}).(logrus.Formatter); ok {
		l.opts.Formatter = formatter
	}
	if hs, ok := l.opts.Context.Value(hooksKey{}).(logrus.LevelHooks); ok {
		l.opts.Hooks = hs
	}
	if caller, ok := l.opts.Context.Value(reportCallerKey{}).(bool); ok && caller {
		l.opts.ReportCaller = caller
	}
	if exitFunction, ok := l.opts.Context.Value(exitKey{}).(func(int)); ok {
		l.opts.ExitFunc = exitFunction
	}

	switch ll := l.opts.Context.Value(loggerKey{}).(type) {
	case *logrus.Logger:
		// overwrite default options
		l.opts.Level = logrusToLoggerLevel(ll.GetLevel())
		l.opts.Out = ll.Out
		l.opts.Formatter = ll.Formatter
		l.opts.Hooks = ll.Hooks
		l.opts.ReportCaller = ll.ReportCaller
		l.opts.ExitFunc = ll.ExitFunc
		l.Logger = ll
	case *logrus.Entry:
		// overwrite default options
		el := ll.Logger
		l.opts.Level = logrusToLoggerLevel(el.GetLevel())
		l.opts.Out = el.Out
		l.opts.Formatter = el.Formatter
		l.opts.Hooks = el.Hooks
		l.opts.ReportCaller = el.ReportCaller
		l.opts.ExitFunc = el.ExitFunc
		l.Logger = ll
	case nil:
		log := logrus.New() // defaults
		log.SetLevel(loggerToLogrusLevel(l.opts.Level))
		log.SetOutput(l.opts.Out)
		log.SetFormatter(l.opts.Formatter)
		log.ReplaceHooks(l.opts.Hooks)
		log.SetReportCaller(l.opts.ReportCaller)
		log.ExitFunc = l.opts.ExitFunc
		l.Logger = log
	default:
		return fmt.Errorf("invalid logrus type: %T", ll)
	}

	return nil
}

func (l *logrusLogger) V(level logger.Level) bool {
	switch ll := l.Logger.(type) {
	case *logrus.Logger:
		return ll.IsLevelEnabled(loggerToLogrusLevel(level))
	case *logrus.Entry:
		return ll.Logger.IsLevelEnabled(loggerToLogrusLevel(level))
	}
	return true
}

func (l *logrusLogger) String() string {
	return "logrus"
}

func (l *logrusLogger) Fields(fields ...interface{}) logger.Logger {
	flds := make(map[string]interface{}, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		flds[fields[i].(string)] = fields[i+1]
	}
	return &logrusLogger{l.Logger.WithFields(flds), l.opts}
}

func (l *logrusLogger) Trace(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.TraceLevel, args...)
}

func (l *logrusLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.TraceLevel, format, args...)
}

func (l *logrusLogger) Warn(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.WarnLevel, args...)
}

func (l *logrusLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.WarnLevel, format, args...)
}

func (l *logrusLogger) Info(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.InfoLevel, args...)
}

func (l *logrusLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.InfoLevel, format, args...)
}

func (l *logrusLogger) Error(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.ErrorLevel, args...)
}

func (l *logrusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.ErrorLevel, format, args...)
}

func (l *logrusLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.FatalLevel, args...)
}

func (l *logrusLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.FatalLevel, format, args...)
}

func (l *logrusLogger) Debug(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.DebugLevel, args...)
}

func (l *logrusLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.DebugLevel, format, args...)
}

func (l *logrusLogger) Log(ctx context.Context, level logger.Level, args ...interface{}) {
	if !l.V(level) {
		return
	}

	l.Logger.Log(loggerToLogrusLevel(level), args...)
}

func (l *logrusLogger) Logf(ctx context.Context, level logger.Level, format string, args ...interface{}) {
	if !l.V(level) {
		return
	}

	l.Logger.Logf(loggerToLogrusLevel(level), format, args...)
}

func (l *logrusLogger) Options() logger.Options {
	// FIXME: How to return full opts?
	return l.opts.Options
}

// New builds a new logger based on options
func NewLogger(opts ...logger.Option) logger.Logger {
	options := Options{
		Options:      logger.NewOptions(opts...),
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: false,
		ExitFunc:     os.Exit,
	}
	l := &logrusLogger{opts: options}
	return l
}

func loggerToLogrusLevel(level logger.Level) logrus.Level {
	switch level {
	case logger.TraceLevel:
		return logrus.TraceLevel
	case logger.DebugLevel:
		return logrus.DebugLevel
	case logger.InfoLevel:
		return logrus.InfoLevel
	case logger.WarnLevel:
		return logrus.WarnLevel
	case logger.ErrorLevel:
		return logrus.ErrorLevel
	case logger.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func logrusToLoggerLevel(level logrus.Level) logger.Level {
	switch level {
	case logrus.TraceLevel:
		return logger.TraceLevel
	case logrus.DebugLevel:
		return logger.DebugLevel
	case logrus.InfoLevel:
		return logger.InfoLevel
	case logrus.WarnLevel:
		return logger.WarnLevel
	case logrus.ErrorLevel:
		return logger.ErrorLevel
	case logrus.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}
