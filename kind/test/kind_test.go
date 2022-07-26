package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"
)

func TestKind(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(true, true)

	logger := NewTestLogger(t, "Kind")
	provider := cluster.NewProvider(cluster.ProviderWithLogger(logger))

	// options := &cluster.CreateOption{}
	err := provider.Create("go-test-1", cluster.CreateWithKubeconfigPath("/tmp/kube-config-kind"))
	t.Log("provider.Create - err=", err)

	err = provider.Delete("go-test-1", "/tmp/kube-config-kind")
	t.Log("Delete - err=", err)
}

type Logger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	WithPrefix(string) Logger
	Write(p []byte) (n int, err error)
	Flush()
}

type kindLogger struct {
	l Logger
}

// TestLogger implements the Logger interface to be compatible with the go test operator's
// output buffering (without this, the use of Parallel tests combined with subtests causes test
// output to be mixed).
type TestLogger struct {
	prefix string
	test   *testing.T
	buffer []byte
}

func NewTestLogger(test *testing.T, prefix string) *TestLogger {
	return &TestLogger{
		prefix: prefix,
		test:   test,
		buffer: []byte{},
	}
}

// Log logs the provided arguments with the logger's prefix. See testing.Log for more details.
func (t *TestLogger) Log(args ...interface{}) {
	args = append([]interface{}{
		fmt.Sprintf("%s | %s |", time.Now().Format("15:04:05"), t.prefix),
	}, args...)
	t.test.Log(args...)
}

// Logf logs the provided arguments with the logger's prefix. See testing.Logf for more details.
func (t *TestLogger) Logf(format string, args ...interface{}) {
	t.Log(fmt.Sprintf(format, args...))
}

// WithPrefix returns a new TestLogger with the provided prefix appended to the current prefix.
func (t *TestLogger) WithPrefix(prefix string) Logger {
	return NewTestLogger(t.test, fmt.Sprintf("%s/%s", t.prefix, prefix))
}

// Write implements the io.Writer interface.
// Logs each line written to it, buffers incomplete lines until the next Write() call.
func (t *TestLogger) Write(p []byte) (n int, err error) {
	t.buffer = append(t.buffer, p...)

	splitBuf := bytes.Split(t.buffer, []byte{'\n'})
	t.buffer = splitBuf[len(splitBuf)-1]

	for _, line := range splitBuf[:len(splitBuf)-1] {
		t.Log(string(line))
	}

	return len(p), nil
}

func (t *TestLogger) Flush() {
	if len(t.buffer) != 0 {
		t.Log(string(t.buffer))
		t.buffer = []byte{}
	}
}

// Warn is part of the log.Logger interface
func (l *TestLogger) Warn(message string) {
	l.Log(message)
}

// Warnf is part of the log.Logger interface
func (l *TestLogger) Warnf(format string, args ...interface{}) {
	l.Logf(format, args...)
}

// Error is part of the log.Logger interface
func (l *TestLogger) Error(message string) {
	l.Log(message)
}

// Errorf is part of the log.Logger interface
func (l *TestLogger) Errorf(format string, args ...interface{}) {
	l.Logf(format, args...)
}

// V is part of the log.Logger interface
func (l *TestLogger) V(level log.Level) log.InfoLogger {
	return infoLogger{
		logger:  l,
		level:   level,
		enabled: true,
	}
}

// debug is like print but with a debug log header
func (l *TestLogger) debug(message string) {
	l.Log(message)
}

// debugf is like printf but with a debug log header
func (l *TestLogger) debugf(format string, args ...interface{}) {
	l.Logf(format, args...)
}

// infoLogger implements log.InfoLogger for Logger
type infoLogger struct {
	logger  *TestLogger
	level   log.Level
	enabled bool
}

// Enabled is part of the log.InfoLogger interface
func (i infoLogger) Enabled() bool {
	return i.enabled
}

// Info is part of the log.InfoLogger interface
func (i infoLogger) Info(message string) {
	if !i.enabled {
		return
	}
	// for > 0, we are writing debug messages, include extra info
	if i.level > 0 {
		i.logger.Logf(message)
	} else {
		i.logger.Logf(message)
	}
}

// Infof is part of the log.InfoLogger interface
func (i infoLogger) Infof(format string, args ...interface{}) {
	if !i.enabled {
		return
	}
	// for > 0, we are writing debug messages, include extra info
	if i.level > 0 {
		i.logger.Logf(format, args...)
	} else {
		i.logger.Logf(format, args...)
	}
}

//
