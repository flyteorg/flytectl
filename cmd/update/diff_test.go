package update

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalToYamlStringRespectsJsonFieldAnnotations(t *testing.T) {
	type T struct {
		FieldIncluded1 int    `json:"fieldIncluded1"`
		FieldIncluded2 string `json:"fieldIncluded2"`
		FieldOmitted   string `json:"fieldOmitted,omitempty"`
	}
	value := T{}

	result, err := marshalToYamlString(value)

	assert.Nil(t, err)
	assert.Equal(t, "fieldIncluded1: 0\nfieldIncluded2: \"\"\n", result)
}

func TestDiffStringsReturnsAUnifiedDiff(t *testing.T) {
	s1 := "abc\ndef\nghi"
	s2 := "aaa\ndef\nghi"

	patch := diffStrings(s1, s2)

	// TODO: Not a unified diff, really - an approximation using diff-match-patch. Need an actual unified diff implementation.
	assert.Equal(t, "@@ -1,8 +1,8 @@\n-abc\n\n+aaa\n\n def\n\n", patch)
}

func TestDiffAsYamlReturnsAUnifiedDiffOfObjectsMarshalledAsYAML(t *testing.T) {
	type T struct {
		F1 int    `json:"f1"`
		F2 string `json:"f2"`
		F3 string `json:"f3,omitempty"`
	}
	object1 := T{F1: 5, F2: "apple"}
	object2 := T{F1: 10, F2: "apple", F3: "banana"}

	patch, err := diffAsYaml(object1, object2)

	assert.Nil(t, err)
	// TODO: Not a unified diff, really - an approximation using diff-match-patch. Need an actual unified diff implementation.
	assert.Equal(t, "@@ -1,10 +1,11 @@\n-f1: 5\n\n+f1: 10\n\n f2: \n@@ -10,8 +10,19 @@\n : apple\n\n+f3: banana\n\n", patch)
}
