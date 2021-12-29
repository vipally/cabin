package logger_test

import (
	"errors"
	"testing"

	"github.com/quanxiang-cloud/cabin/logger"
)

func TestLogger(t *testing.T) {
	cfg := &logger.Config{
		Level: logger.DebugLevel.Int(),
	}
	logger.Default = logger.New(cfg)

	namedLog := logger.Default.WithName("named")
	log := logger.Default

	namedLog.Infof("info %s", "foo")
	namedLog.PutError(nil, "no error")
	namedLog.PutError(errors.New("err"), "")
	log.Infof("info %s", "foo")
	log.Info("info")
	log.Debug("debug")
	log.Warn("warn")
	log.Error("err")
	log.Panic("panic")
	log.Fatal("fatal")

}
