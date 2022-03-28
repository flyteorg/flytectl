package testutils

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/flyteorg/flyteidl/clients/go/admin"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	extMocks "github.com/flyteorg/flytectl/pkg/ext/mocks"
	"github.com/stretchr/testify/assert"
)

const projectValue = "dummyProject"
const domainValue = "dummyDomain"
const output = "json"

var (
	reader        *os.File
	writer        *os.File
	Err           error
	Ctx           context.Context
	MockClient    *admin.Clientset
	FetcherExt    *extMocks.AdminFetcherExtInterface
	UpdaterExt    *extMocks.AdminUpdaterExtInterface
	DeleterExt    *extMocks.AdminDeleterExtInterface
	MockOutStream io.Writer
	CmdCtx        cmdCore.CommandContext
	stdOut        *os.File
	stderr        *os.File
)

func Setup() {
	Ctx = context.Background()
	reader, writer, Err = os.Pipe()
	if Err != nil {
		panic(Err)
	}
	stdOut = os.Stdout
	stderr = os.Stderr
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	MockClient = admin.InitializeMockClientset()
	FetcherExt = new(extMocks.AdminFetcherExtInterface)
	UpdaterExt = new(extMocks.AdminUpdaterExtInterface)
	DeleterExt = new(extMocks.AdminDeleterExtInterface)
	FetcherExt.OnAdminServiceClient().Return(MockClient.AdminClient())
	UpdaterExt.OnAdminServiceClient().Return(MockClient.AdminClient())
	DeleterExt.OnAdminServiceClient().Return(MockClient.AdminClient())
	MockOutStream = writer
	CmdCtx = cmdCore.NewCommandContextWithExt(MockClient, FetcherExt, UpdaterExt, DeleterExt, MockOutStream)
	config.GetConfig().Project = projectValue
	config.GetConfig().Domain = domainValue
	config.GetConfig().Output = output
}

// TearDownAndVerify TODO: Change this to verify log lines from context
func TearDownAndVerify(t *testing.T, expectedLog string) {
	writer.Close()
	os.Stdout = stdOut
	os.Stderr = stderr
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err == nil {
		assert.Equal(t, sanitizeString(expectedLog), sanitizeString(buf.String()))
	}
}

func sanitizeString(str string) string {
	// Not the most comprehensive ANSI pattern, but this should capture common color operations such as \x1b[107;0m and \x1b[0m. Expand if needed (insert regex 2 problems joke here).
	ansiRegex := regexp.MustCompile("\u001B\\[[\\d+\\;]*\\d+m")
	return ansiRegex.ReplaceAllString(strings.Trim(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, "\n", ""), "\t", ""), "", ""), " \t"), "")
}
