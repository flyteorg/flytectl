package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

func TestSetupConfig(t *testing.T) {
	assert.Nil(t, initFlytectlConfig(strings.NewReader("Yes")))
	assert.Nil(t, initFlytectlConfig(strings.NewReader("Yes")))
	assert.Nil(t, initFlytectlConfig(strings.NewReader("No")))
}

func TestConfigCmd(t *testing.T) {
	assert.Nil(t, configCmd.RunE(&cobra.Command{}, []string{"init"}))
}
