package cmdcore

import (
	"fmt"
	interfaces2 "github.com/flyteorg/flytectl/cmd/get/interfaces"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
	"io"
)

// CmdCtxBuilder is used to build the CmdCtx.
type CmdCtxBuilder struct {
	adminClient service.AdminServiceClient
	in          io.Reader
	out         io.Writer
	fetcher     interfaces2.Fetcher
}

// CmdContextBuilder is constructor function to be used by the clients in interacting with the builder
func CmdContextBuilder() *CmdCtxBuilder {
	return &CmdCtxBuilder{}
}

// WithAdminClient provides the adminClient to be used for constructing the command context
func (cb *CmdCtxBuilder) WithAdminClient(adminClient service.AdminServiceClient) *CmdCtxBuilder {
	cb.adminClient = adminClient
	return cb
}

// WithReader allows pluggable reader for command context
func (cb *CmdCtxBuilder) WithReader(in io.Reader) *CmdCtxBuilder {
	cb.in = in
	return cb
}

// WithWriter allows pluggable writer for command context
func (cb *CmdCtxBuilder) WithWriter(out io.Writer) *CmdCtxBuilder {
	cb.out = out
	return cb
}

// WithFetcher allows to plugin the fetcher if required to get data from admin
func (cb *CmdCtxBuilder) WithFetcher(fetcher interfaces2.Fetcher) *CmdCtxBuilder {
	cb.fetcher = fetcher
	return cb
}

// Build the commandcontext using the current state of the CmdCtxBuilder
func (cb *CmdCtxBuilder) Build() (CommandContext, error) {
	if cb.adminClient == nil {
		return CommandContext{}, fmt.Errorf("admin client must be set for builder")
	}

	if cb.out == nil {
		return CommandContext{}, fmt.Errorf("writer must be set for builder")
	}

	return NewCmdCtxUsingBuilder(cb.adminClient, cb.in, cb.out, cb.fetcher), nil
}
