package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteIntoFile(t *testing.T) {
	t.Run("Successfully write into a file", func(t *testing.T) {
		err := WriteIntoFile([]byte(""), "version.yaml")
		assert.Nil(t, err)
	})
	t.Run("Error in writing file", func(t *testing.T) {
		err := WriteIntoFile([]byte(""), "version.yaml")
		assert.Nil(t, err)
	})
}

func TestSetupFlyteDir(t *testing.T) {
	assert.Nil(t, SetupFlyteDir())
}

func TestPrintSandboxMessage(t *testing.T) {
	t.Run("Print Sandbox Message", func(t *testing.T) {
		PrintSandboxMessage()
	})
}

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
