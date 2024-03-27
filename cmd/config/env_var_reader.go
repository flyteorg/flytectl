package config

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/flyteorg/flyte/flyteidl/clients/go/admin"
	"github.com/flyteorg/flyte/flytestdlib/config"
)

type FuncType func() error

var funcMap = map[string]FuncType{}

func init() {
	funcMap["FLYTE_ADMIN_ENDPOINT"] = getAdminEndpoint
	// TODO add more env vars if needed
}

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

func getAdminEndpoint() error {
	ctx := context.Background()
	cfg := admin.GetConfig(ctx)
	if len(os.Getenv("FLYTE_ADMIN_ENDPOINT")) > 0 {
		envEndpoint, err := url.Parse(os.Getenv("FLYTE_ADMIN_ENDPOINT"))
		if err != nil {
			return fmt.Errorf("error parsing env var flyte_admin_endpoint: %v", err)
		}
		cfg.Endpoint = config.URL{URL: *envEndpoint}
		admin.SetConfig(cfg)
	}
	return nil
}
