package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmdIntegration(t *testing.T) {
	rootCmd, err := newRootCmd()
	assert.Nil(t, err)
	assert.NotNil(t, rootCmd)
}
