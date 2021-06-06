package serialize

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/flyteorg/flytectl/cmd/config"
	sconfig "github.com/flyteorg/flytectl/cmd/config/subcommand/serialize"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/clients/go/admin"
)

const (
	serializeWorkflowShort = "Serialize flyte workflow"
	serializeWorkflowLong  = `
Serialize workflow
::

 bin/flytectl serialize workflow  -p flytesnacks -d development --image="flyteorg/flytecookbook:core-1d1631120a2e3505a918aa536f6f7a71b52147a3" --service-account="default" --version="v22" --output-dir="flyteorg/flytesnacks/cookbook/core/_pb_output" --output-dir-prefix="s3://my-s3-bucket/raw-data" --flyte-aws-endpoint="http://localhost:30084/" --flyte-aws-key="minio" --flyte-aws-secret="miniostorage" --command="serialize"

Usage
`
)

func serializeWorkflowFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	if len(sconfig.DefaultFilesConfig.Image) > 0 {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return err
		}
		r, err := cli.ImagePull(ctx, sconfig.DefaultFilesConfig.Image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		if _, err := io.Copy(os.Stdout, r); err != nil {
			return err
		}
		c := admin.GetConfig(ctx)

		environment := []string{
			fmt.Sprintf("REGISTRY=%s", ""),
			fmt.Sprintf("FLYTE_HOST=%s", sconfig.DefaultFilesConfig.Registry),
			fmt.Sprintf("PROJECT=%s", config.GetConfig().Project),
			fmt.Sprintf("SERVICE_ACCOUNT=%s", sconfig.DefaultFilesConfig.ServiceAccount),
			fmt.Sprintf("VERSION=%s", sconfig.DefaultFilesConfig.Version),
			fmt.Sprintf("OUTPUT_DATA_PREFIX=%s", sconfig.DefaultFilesConfig.OutputDirprefix),
			fmt.Sprintf("FLYTE_AWS_ENDPOINT=%s", sconfig.DefaultFilesConfig.FlyteAwsEndpoint),
			fmt.Sprintf("FLYTE_AWS_ACCESS_KEY_ID=%s", sconfig.DefaultFilesConfig.FlyteAwsKey),
			fmt.Sprintf("FLYTE_AWS_SECRET_ACCESS_KEY=%s", sconfig.DefaultFilesConfig.FlyteAwsSecret),
			fmt.Sprintf("INSECURE_FLAG=%v", c.UseInsecureConnection),
			fmt.Sprintf("FLYTE_AWS_SECRET_ACCESS_KEY=%s", sconfig.DefaultFilesConfig.FlyteAwsSecret),
			fmt.Sprintf("ADDL_DISTRIBUTION_DIR=%s", "s3://my-s3-bucket/fast/"),
		}

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Env:   environment,
			Image: sconfig.DefaultFilesConfig.Image,
			Tty:   false,
			Entrypoint: strslice.StrSlice{
				"make", sconfig.DefaultFilesConfig.Command,
			},
		}, &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: sconfig.DefaultFilesConfig.OutputDir,
					Target: "/tmp/output",
				},
			},
			Privileged: true,
		}, nil,
			nil, "flytecookbook")
		if err != nil {
			return err
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}

		reader, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: false,
			Follow:     true,
		})
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		return nil
	}
	return fmt.Errorf("Please specify --image flag for docker imag")
}
