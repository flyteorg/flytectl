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

type FuncType func() error

var funcMap = map[string]FuncType{flyteAdminEndpoint: updateAdminEndpoint}

func UpdateConfigWithEnvVar() error {
	for envVar, f := range funcMap {
		if os.Getenv(envVar) != "" {
			if err := f(); err != nil {
				return fmt.Errorf("error update config with env var: %v", err)
			}
		}
	}
	return nil
}

func updateAdminEndpoint() error {
	ctx := context.Background()
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
