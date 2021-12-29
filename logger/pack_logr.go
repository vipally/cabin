package logger

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var _ Logger = (*packLogr)(nil)

// packLogr adapt logr.Loger to Logger
type packLogr struct {
	l            logr.Logger
	minLevel     Level
	keyAndValues []interface{}
}

func newPackLogr(cfg *Config) *packLogr {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	p := &packLogr{
		l:            zapr.NewLoggerWithOptions(zapLogger).WithCallDepth(1),
		minLevel:     DebugLevel,
		keyAndValues: nil,
	}
	return p.init(cfg)
}

func (p *packLogr) init(cfg *Config) *packLogr {
	if cfg != nil {
		p.minLevel = Level(cfg.Level)
	}
	return p
}

func (p *packLogr) log(lv Level, format string, fmtArgs []interface{}, context []interface{}) {
	if lv < DPanicLevel && lv < p.minLevel { //ignore level with too low priority
		//return
	}
	msg := format
	switch {
	case msg == "" && len(fmtArgs) > 0:
		msg = fmt.Sprint(fmtArgs...)
	case msg != "" && len(fmtArgs) > 0:
		msg = fmt.Sprintf(format, fmtArgs...)
	}

	args := append(p.keyAndValues, context...)
	if cap(p.keyAndValues) < len(args) {
		p.keyAndValues = args[:len(p.keyAndValues)]
	}
	v := p.l.V(lv.index())
	fmt.Printf("%d %#v\n", lv, v)
	v.Info(msg, args...)
}

// WithValues returns a new Logger with additional key/value pairs.
func (p *packLogr) WithValues(keysAndValues ...interface{}) Logger {
	n := &packLogr{
		l:            p.l.WithValues(keysAndValues...),
		minLevel:     p.minLevel,
		keyAndValues: append([]interface{}(nil), p.keyAndValues...),
	}
	return n
}

// WithName returns a new Logger with the specified name appended.
func (p *packLogr) WithName(name string) Logger {
	n := &packLogr{
		l:            p.l.WithName(name),
		minLevel:     p.minLevel,
		keyAndValues: append([]interface{}(nil), p.keyAndValues...),
	}
	return n
}

// PutError write log with error
func (p *packLogr) PutError(err error, msg string, keysAndValues ...interface{}) {
	p.l.Error(err, msg, keysAndValues...)
}

// Debug uses fmt.Sprint to construct and log a message.
func (p *packLogr) Debug(args ...interface{}) {
	p.log(DebugLevel, "", args, nil)
}

// Info uses fmt.Sprint to construct and log a message.
func (p *packLogr) Info(args ...interface{}) {
	p.log(InfoLevel, "", args, nil)
}

// Warn uses fmt.Sprint to construct and log a message.
func (p *packLogr) Warn(args ...interface{}) {
	p.log(WarnLevel, "", args, nil)
}

// Error uses fmt.Sprint to construct and log a message.
func (p *packLogr) Error(args ...interface{}) {
	p.log(ErrorLevel, "", args, nil)
}

// Panic uses fmt.Sprint to construct and log a message, then panicl.
func (p *packLogr) Panic(args ...interface{}) {
	p.log(PanicLevel, "", args, nil)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls ol.Exit.
func (p *packLogr) Fatal(args ...interface{}) {
	p.log(FatalLevel, "", args, nil)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (p *packLogr) Debugf(template string, args ...interface{}) {
	p.log(DebugLevel, template, args, nil)
}

// Infof uses fmt.Sprintf to log a templated message.
func (p *packLogr) Infof(template string, args ...interface{}) {
	p.log(InfoLevel, template, args, nil)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (p *packLogr) Warnf(template string, args ...interface{}) {
	p.log(WarnLevel, template, args, nil)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (p *packLogr) Errorf(template string, args ...interface{}) {
	p.log(ErrorLevel, template, args, nil)
}

// Panicf uses fmt.Sprintf to log a templated message, then panicl.
func (p *packLogr) Panicf(template string, args ...interface{}) {
	p.log(PanicLevel, template, args, nil)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls ol.Exit.
func (p *packLogr) Fatalf(template string, args ...interface{}) {
	p.log(FatalLevel, template, args, nil)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  l.With(keysAndValues).Debug(msg)
func (p *packLogr) Debugw(msg string, keysAndValues ...interface{}) {
	p.log(DebugLevel, msg, nil, keysAndValues)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (p *packLogr) Infow(msg string, keysAndValues ...interface{}) {
	p.log(InfoLevel, msg, nil, keysAndValues)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (p *packLogr) Warnw(msg string, keysAndValues ...interface{}) {
	p.log(WarnLevel, msg, nil, keysAndValues)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (p *packLogr) Errorw(msg string, keysAndValues ...interface{}) {
	p.log(ErrorLevel, msg, nil, keysAndValues)
}

// Panicw logs a message with some additional context, then panicl. The
// variadic key-value pairs are treated as they are in With.
func (p *packLogr) Panicw(msg string, keysAndValues ...interface{}) {
	p.log(PanicLevel, msg, nil, keysAndValues)
}

// Fatalw logs a message with some additional context, then calls ol.Exit. The
// variadic key-value pairs are treated as they are in With.
func (p *packLogr) Fatalw(msg string, keysAndValues ...interface{}) {
	p.log(FatalLevel, msg, nil, keysAndValues)
}

// Sync flushes any buffered log entries.
func (p *packLogr) Sync() error {
	// NOTE: do nothing for logr
	return nil
}
