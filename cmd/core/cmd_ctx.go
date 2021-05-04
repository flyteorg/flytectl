package cmdcore

import (
	"io"

	"github.com/flyteorg/flytectl/cmd/get/interfaces"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
)

type CommandContext struct {
	adminClient service.AdminServiceClient
	in          io.Reader
	out         io.Writer
	fetcher    interfaces.Fetcher
}
// NewCommandContext Deprecated. Use the NewCmdCtxUsingBuilder
func NewCommandContext(adminClient service.AdminServiceClient, out io.Writer) CommandContext {
	return CommandContext{adminClient: adminClient, out: out}
}

func NewCmdCtxUsingBuilder(adminClient service.AdminServiceClient, in io.Reader, out io.Writer,
	fetcher interfaces.Fetcher) CommandContext {
	return CommandContext{adminClient: adminClient, in: in , out: out, fetcher: fetcher}
}

func (c CommandContext) AdminClient() service.AdminServiceClient {
	return c.adminClient
}

func (c CommandContext) OutputPipe() io.Writer {
	return c.out
}

func (c CommandContext) InputPipe() io.Reader {
	return c.in
}

func (c CommandContext) Fetcher() interfaces.Fetcher {
	return c.fetcher
}
