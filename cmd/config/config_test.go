package config

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/flyteorg/flyte/flyteidl/clients/go/admin"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/stretchr/testify/assert"
)

func TestOutputFormat(t *testing.T) {
	c := &Config{
		Output: "json",
	}
	result, err := c.OutputFormat()
	assert.Nil(t, err)
	assert.Equal(t, printer.OutputFormat(1), result)
}

func TestInvalidOutputFormat(t *testing.T) {
	c := &Config{
		Output: "test",
	}
	var result printer.OutputFormat
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, printer.OutputFormat(0), result)
			assert.NotNil(t, r)
		}
	}()
	result = c.MustOutputFormat()
}

func TestUpdateConfigWithEnvVar(t *testing.T) {
	originalValue := os.Getenv("FLYTE_ADMIN_ENDPOINT")
	defer os.Setenv("FLYTE_ADMIN_ENDPOINT", originalValue)

	dummyUrl := "dns://dummyHost"
	os.Setenv("FLYTE_ADMIN_ENDPOINT", dummyUrl)

	parsedDummyUrl, _ := url.Parse(dummyUrl)

	adminCfg := admin.GetConfig(context.Background())

	assert.NotEqual(t, adminCfg.Endpoint.URL, *parsedDummyUrl)
	UpdateConfigWithEnvVar()
	assert.Equal(t, adminCfg.Endpoint.URL, *parsedDummyUrl)
}
