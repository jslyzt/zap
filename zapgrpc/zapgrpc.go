package zapgrpc

import "github.com/jslyzt/zap"

// An Option overrides a Logger's default configuration.
type Option interface {
	apply(*Logger)
}

type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

// WithDebug configures a Logger to print at zap's DebugLevel instead of
// InfoLevel.
func WithDebug() Option {
	return optionFunc(func(logger *Logger) {
		logger.print = (*zap.SugaredLogger).Debug
		logger.printf = (*zap.SugaredLogger).Debugf
	})
}

// NewLogger returns a new Logger.
//
// By default, Loggers print at zap's InfoLevel.
func NewLogger(l *zap.Logger, options ...Option) *Logger {
	logger := &Logger{
		log:    l.Sugar(),
		fatal:  (*zap.SugaredLogger).Fatal,
		fatalf: (*zap.SugaredLogger).Fatalf,
		print:  (*zap.SugaredLogger).Info,
		printf: (*zap.SugaredLogger).Infof,
	}
	for _, option := range options {
		option.apply(logger)
	}
	return logger
}

// Logger adapts zap's Logger to be compatible with grpclog.Logger.
type Logger struct {
	log    *zap.SugaredLogger
	fatal  func(*zap.SugaredLogger, ...interface{})
	fatalf func(*zap.SugaredLogger, string, ...interface{})
	print  func(*zap.SugaredLogger, ...interface{})
	printf func(*zap.SugaredLogger, string, ...interface{})
}

// Fatal implements grpclog.Logger.
func (l *Logger) Fatal(args ...interface{}) {
	l.fatal(l.log, args...)
}

// Fatalf implements grpclog.Logger.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.fatalf(l.log, format, args...)
}

// Fatalln implements grpclog.Logger.
func (l *Logger) Fatalln(args ...interface{}) {
	l.fatal(l.log, args...)
}

// Print implements grpclog.Logger.
func (l *Logger) Print(args ...interface{}) {
	l.print(l.log, args...)
}

// Printf implements grpclog.Logger.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.printf(l.log, format, args...)
}

// Println implements grpclog.Logger.
func (l *Logger) Println(args ...interface{}) {
	l.print(l.log, args...)
}
