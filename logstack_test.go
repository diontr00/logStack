package logStack

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestLoggerSuite struct {
	suite.Suite
	logger *Logger
	buf    *bytes.Buffer
}

func (t *TestLoggerSuite) SetupTest() {
	t.buf = new(bytes.Buffer)
	t.logger = NewLogger(t.buf, InfoLevel)
}

func (t *TestLoggerSuite) TestInfoLevel() {
	t.buf.Reset()
	msg := "testing info"
	fields := []LogField{String("[INFO]", "diontr"), Int("attemp", 1)}
	infoWant := fmt.Sprintf(
		`{"level":"info","ts":"%s","msg":"%s","[INFO]":"diontr","attemp":1}`,
		time.Now().Format(time.RFC1123), msg,
	)

	t.logger.Info(msg, fields...)
	assert.Equal(t.T(), infoWant, strings.TrimSuffix(t.buf.String(), "\n"))
}

func (t *TestLoggerSuite) TestEmpty() {
	t.buf.Reset()
	msg := "testing debug"
	fields := []LogField{String("[DEBUG]", "diontr"), Int("attemp", 1)}
	debugWant := ""
	t.logger.Debug(msg, fields...)
	assert.Equal(t.T(), debugWant, strings.TrimSuffix(t.buf.String(), "\n"))

}

func (t *TestLoggerSuite) TestMultiLogger() {
	var errBuf bytes.Buffer
	var infoBuf bytes.Buffer

	multiOps := []MultiOption{
		{
			W: &infoBuf,
			Level: func(lvl Loglevel) bool {
				return lvl <= InfoLevel
			},
		},
		{
			W: &errBuf,
			Level: func(lvl Loglevel) bool {
				return lvl >= ErrorLevel
			},
		},
	}

	multiLogger := NewMultiLoger(multiOps)

	ResetDefault(multiLogger)

	infomsg := "testing info"
	infofields := []LogField{String("[INFO]", "diontr"), Int("attemp", 1)}
	infoWant := fmt.Sprintf(
		`{"level":"info","ts":"%s","msg":"%s","[INFO]":"diontr","attemp":1}`,
		time.Now().Format(time.RFC1123), infomsg,
	)

	errmsg := "testing error"
	errorfields := []LogField{String("[ERROR]", "diontr"), Int("attemp", 1)}
	errorWant := fmt.Sprintf(
		`{"level":"error","ts":"%s","msg":"%s","[ERROR]":"diontr","attemp":1}`,
		time.Now().Format(time.RFC1123), errmsg,
	)

	multiLogger.Info(infomsg, infofields...)
	multiLogger.Error(errmsg, errorfields...)

	assert.Equal(t.T(), infoWant, strings.TrimSuffix(infoBuf.String(), "\n"))
	assert.Equal(t.T(), errorWant, strings.TrimSuffix(errBuf.String(), "\n"))

}

func TestLogger(t *testing.T) {
	suite.Run(t, new(TestLoggerSuite))

}
