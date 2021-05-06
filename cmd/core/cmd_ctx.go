package cmdcore

import (
	"io"

	"github.com/flyteorg/flytectl/pkg/ext"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
)

type CommandContext struct {
	adminClient           service.AdminServiceClient
	adminClientFetcherExt ext.AdminFetcherExtInterface
	adminClientUpdateExt  ext.AdminUpdaterExtInterface
	adminClientDeleteExt  ext.AdminDeleterExtInterface
	in                    io.Reader
	out                   io.Writer
}

func NewCommandContext(adminClient service.AdminServiceClient, out io.Writer) CommandContext {
	return CommandContext{adminClient: adminClient, out: out,
		adminClientFetcherExt: &ext.AdminFetcherExtClient{AdminClient: adminClient},
		adminClientUpdateExt:  &ext.AdminUpdaterExtClient{AdminClient: adminClient},
		adminClientDeleteExt:  &ext.AdminDeleterExtClient{AdminClient: adminClient}}
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

func (c CommandContext) AdminFetcherExt() ext.AdminFetcherExtInterface {
	return c.adminClientFetcherExt
}

func (c CommandContext) AdminUpdaterExt() ext.AdminUpdaterExtInterface {
	return c.adminClientUpdateExt
}

func (c CommandContext) AdminDeleterExt() ext.AdminDeleterExtInterface {
	return c.adminClientDeleteExt
}
