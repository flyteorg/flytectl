package cmdcore

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/flyteorg/flyte/flyteidl/clients/go/admin"
	stdConfig "github.com/flyteorg/flyte/flytestdlib/config"
	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flytectl/pkg/pkce"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type PFlagProvider interface {
	GetPFlagSet(prefix string) *pflag.FlagSet
}

type CommandEntry struct {
	ProjectDomainNotRequired bool
	CmdFunc                  CommandFunc
	Aliases                  []string
	Short                    string
	Long                     string
	PFlagProvider            PFlagProvider
	DisableFlyteClient       bool
}

func AddCommands(rootCmd *cobra.Command, cmdFuncs map[string]CommandEntry) {
	for resource, cmdEntry := range cmdFuncs {
		cmd := &cobra.Command{
			Use:          resource,
			Short:        cmdEntry.Short,
			Long:         cmdEntry.Long,
			Aliases:      cmdEntry.Aliases,
			RunE:         generateCommandFunc(cmdEntry),
			SilenceUsage: true,
		}

		if cmdEntry.PFlagProvider != nil {
			cmd.Flags().AddFlagSet(cmdEntry.PFlagProvider.GetPFlagSet(""))
		}

		rootCmd.AddCommand(cmd)
	}
}

func generateCommandFunc(cmdEntry CommandEntry) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if !cmdEntry.ProjectDomainNotRequired {
			if config.GetConfig().Project == "" {
				return fmt.Errorf("project and domain are required parameters")
			}
			if config.GetConfig().Domain == "" {
				return fmt.Errorf("project and domain are required parameters")
			}
		}
		if _, err := config.GetConfig().OutputFormat(); err != nil {
			return err
		}

		adminCfg := admin.GetConfig(ctx)

		if len(os.Getenv("FLTYE_ADMIN_ENDPOINT")) > 0 {
			envEndpoint, err := url.Parse(os.Getenv("FLTYE_ADMIN_ENDPOINT"))
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return err
			}
			adminCfg.Endpoint = stdConfig.URL{
				URL: *envEndpoint,
			}
		}

		if len(adminCfg.Endpoint.String()) == 0 {
			return cmdEntry.CmdFunc(ctx, args, CommandContext{})
		}

		cmdCtx := NewCommandContextNoClient(cmd.OutOrStdout())
		if !cmdEntry.DisableFlyteClient {
			clientSet, err := admin.ClientSetBuilder().WithConfig(admin.GetConfig(ctx)).
				WithTokenCache(pkce.TokenCacheKeyringProvider{
					ServiceUser: fmt.Sprintf("%s:%s", adminCfg.Endpoint.String(), pkce.KeyRingServiceUser),
					ServiceName: pkce.KeyRingServiceName,
				}).Build(ctx)
			if err != nil {
				return err
			}
			cmdCtx = NewCommandContext(clientSet, cmd.OutOrStdout())
		}

		err := cmdEntry.CmdFunc(ctx, args, cmdCtx)
		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.Unavailable || s.Code() == codes.Unauthenticated || s.Code() == codes.Unknown {
					return errors.WithMessage(err,
						fmt.Sprintf("Connection Info: [Endpoint: %s, InsecureConnection?: %v, AuthMode: %v]", adminCfg.Endpoint.String(), adminCfg.UseInsecureConnection, adminCfg.AuthType))
				}
			}
			return err
		}
		return nil
	}
}
