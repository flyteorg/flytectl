//go:build integration
// +build integration

package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendRequest(t *testing.T) {
	t.Run("Successful get request", func(t *testing.T) {
		response, err := SendRequest("GET", "https://github.com", nil)
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("Successful get request failed", func(t *testing.T) {
		response, err := SendRequest("GET", "htp://github.com", nil)
		assert.NotNil(t, err)
		assert.Nil(t, response)
	})
	t.Run("Successful get request failed", func(t *testing.T) {
		response, err := SendRequest("GET", "https://github.com/evalsocket/flyte/archive/refs/tags/source-code.zip", nil)
		assert.NotNil(t, err)
		assert.Nil(t, response)
	})
}
