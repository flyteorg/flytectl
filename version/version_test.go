package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
)

type dFormat struct {
}

var testdata = map[string]string{
	"app":        "hello",
	"version":    "c00b0b6",
	"build":      "c00b0b6",
	"build-time": "2021-03-25",
}

func (dFormat) Format(e *logrus.Entry) ([]byte, error) {
	return []byte(e.Message), nil
}

func TestLogBuildInformation(t *testing.T) {

	buf := bytes.NewBufferString("")
	logrus.SetFormatter(dFormat{})
	logrus.SetOutput(buf)

	Build = testdata["build"]
	BuildTime = testdata["build-time"]
	Version = testdata["version"]
	b, _ := json.Marshal(testdata)
	LogBuildInformation(testdata["app"])
	assert.Equal(t, buf.String(), fmt.Sprintf(string(b)))
}
