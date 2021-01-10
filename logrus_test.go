package logrus

import (
	"errors"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/unistack-org/micro/v3/logger"
)

func TestName(t *testing.T) {
	l := NewLogger()
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.V(logger.InfoLevel)
	if l.String() != "logrus" {
		t.Errorf("error: name expected 'logrus' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

func TestWithFields(t *testing.T) {
	l := NewLogger(logger.WithOutput(os.Stdout)).Fields(map[string]interface{}{
		"k1": "v1",
		"k2": 123456,
	})
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	logger.DefaultLogger = l

	logger.Info("testing: Info")
	logger.Infof("testing: %s", "Infof")
}

func TestWithError(t *testing.T) {
	l := NewLogger().Fields(map[string]interface{}{"error": errors.New("boom!")})
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	logger.DefaultLogger = l

	logger.Error("testing: error")
}

func TestWithLogger(t *testing.T) {
	// with *logrus.Logger
	l := NewLogger(WithLogger(logrus.StandardLogger())).Fields(map[string]interface{}{
		"k1": "v1",
		"k2": 123456,
	})
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	logger.DefaultLogger = l
	logger.Info("testing: with *logrus.Logger")

	// with *logrus.Entry
	el := NewLogger(WithLogger(logrus.NewEntry(logrus.New()))).Fields(map[string]interface{}{
		"k3": 3.456,
		"k4": true,
	})
	if err := el.Init(); err != nil {
		t.Fatal(err)
	}

	logger.DefaultLogger = el
	logger.Info("testing: with *logrus.Entry")
}

func TestJSON(t *testing.T) {
	logger.DefaultLogger = NewLogger(WithJSONFormatter(&logrus.JSONFormatter{}))

	logger.Infof("test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	logger.DefaultLogger = NewLogger()

	if err := logger.Init(logger.WithLevel(logger.DebugLevel)); err != nil {
		t.Fatal(err)
	}

	logger.Debugf("test show debug: %s", "debug msg")

	if err := logger.Init(logger.WithLevel(logger.InfoLevel)); err != nil {
		t.Fatal(err)
	}

	logger.Debugf("test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	l := NewLogger(ReportCaller())

	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	logger.DefaultLogger = l
	logger.Infof("testing: %s", "WithReportCaller")
}
