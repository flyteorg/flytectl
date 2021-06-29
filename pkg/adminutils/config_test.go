package adminutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	DefaultConfig = &Config{
		MaxRecords: 500,
		BatchSize:  100,
	}
	c := GetConfig()
	assert.Equal(t, DefaultConfig.BatchSize, c.BatchSize)
	assert.Equal(t, DefaultConfig.MaxRecords, c.MaxRecords)
}
