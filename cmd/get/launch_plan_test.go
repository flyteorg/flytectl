package get

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_transformLaunchPlan(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		v, err := json.Marshal(map[string]string{
			"Version":          "id",
			"Name":             "name",
			"Type":             "type",
			"DiscoveryVersion": "discoveryVersion",
		})
		assert.NoError(t, err)
		row, err := transformLaunchPlan(v)
		assert.NoError(t, err)
		typedRow, ok := row.(PrintableLaunchPlan)
		assert.True(t, ok)
		assert.NotNil(t, typedRow)
		assert.Equal(t, "name", typedRow.Name)
	})
}
