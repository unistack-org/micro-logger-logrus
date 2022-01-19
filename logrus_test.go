package logrus

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"go.unistack.org/micro/v3/logger"
)

func TestFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Fields("key", "val").Info(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`key=val`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}
}

func TestOutput(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Infof(context.TODO(), "test logger name: %s", "name")
	if !bytes.Contains(buf.Bytes(), []byte(`test logger name`)) {
		t.Fatalf("log not redirected: %s", buf.Bytes())
	}
}

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
	l := NewLogger(logger.WithOutput(os.Stdout))

	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l = l.Fields("k1", "v1", "k2", 123456)

	logger.DefaultLogger = l

	logger.Info(context.TODO(), "testing: Info")
	logger.Infof(context.TODO(), "testing: %s", "Infof")
}

func TestWithError(t *testing.T) {
	l := NewLogger()
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l = l.Fields("error", errors.New("boom!"))

	logger.DefaultLogger = l

	logger.Error(context.TODO(), "testing: error")
}

func TestWithLogger(t *testing.T) {
	// with *logrus.Logger
	l := NewLogger(WithLogger(logrus.StandardLogger()))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l = l.Fields("k1", "v1", "k2", 123456)

	logger.DefaultLogger = l
	logger.Info(context.TODO(), "testing: with *logrus.Logger")

	// with *logrus.Entry
	el := NewLogger(WithLogger(logrus.NewEntry(logrus.New())))
	if err := el.Init(); err != nil {
		t.Fatal(err)
	}

	l = l.Fields("k3", 3.456, "k4", true)

	logger.DefaultLogger = el
	logger.Info(context.TODO(), "testing: with *logrus.Entry")
}

func TestJSON(t *testing.T) {
	l := NewLogger(WithJSONFormatter(&logrus.JSONFormatter{}))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	logger.DefaultLogger = l
	logger.Infof(context.TODO(), "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	logger.DefaultLogger = NewLogger()

	if err := logger.Init(logger.WithLevel(logger.DebugLevel)); err != nil {
		t.Fatal(err)
	}

	logger.Debugf(context.TODO(), "test show debug: %s", "debug msg")

	if err := logger.Init(logger.WithLevel(logger.InfoLevel)); err != nil {
		t.Fatal(err)
	}

	logger.Debugf(context.TODO(), "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	l := NewLogger(ReportCaller())

	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	logger.DefaultLogger = l
	logger.Infof(context.TODO(), "testing: %s", "WithReportCaller")
}
