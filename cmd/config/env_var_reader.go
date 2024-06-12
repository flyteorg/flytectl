package config

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/flyteorg/flyte/flyteidl/clients/go/admin"
	"github.com/flyteorg/flyte/flytestdlib/config"
)

const flyteAdminEndpoint = "FLYTE_ADMIN_ENDPOINT"

type UpdateFunc func(context.Context) error

var envToUpdateFunc = map[string]UpdateFunc{flyteAdminEndpoint: updateAdminEndpoint}

func UpdateConfigWithEnvVar() error {
	ctx := context.Background()

	for envVar, updateFunc := range envToUpdateFunc {
		if os.Getenv(envVar) != "" {
			if err := updateFunc(ctx); err != nil {
				return fmt.Errorf("error update config with env var: %v", err)
			}
		}
	}
	return nil
}

func updateAdminEndpoint(ctx context.Context) error {
	cfg := admin.GetConfig(ctx)

	if len(os.Getenv(flyteAdminEndpoint)) > 0 {
		envEndpoint, err := url.Parse(os.Getenv(flyteAdminEndpoint))
		if err != nil {
			return fmt.Errorf("error parsing env var %v: %v", flyteAdminEndpoint, err)
		}
		cfg.Endpoint = config.URL{URL: *envEndpoint}
		if err := admin.SetConfig(cfg); err != nil {
			return err
		}
	}
	return nil
}
