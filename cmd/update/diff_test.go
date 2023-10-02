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

	patch := diffStrings("before", "after", s1, s2)

	assert.Equal(t, "--- before\n+++ after\n@@ -1,3 +1,3 @@\n-abc\n+aaa\n def\n ghi\n", patch)
}

func TestDiffAsYamlReturnsAUnifiedDiffOfObjectsMarshalledAsYAML(t *testing.T) {
	type T struct {
		F1 int    `json:"f1"`
		F2 string `json:"f2"`
		F3 string `json:"f3,omitempty"`
	}
	object1 := T{F1: 5, F2: "apple"}
	object2 := T{F1: 10, F2: "apple", F3: "banana"}

	patch, err := DiffAsYaml("before", "after", object1, object2)

	assert.Nil(t, err)
	assert.Equal(t, "--- before\n+++ after\n@@ -1,3 +1,4 @@\n-f1: 5\n+f1: 10\n f2: apple\n+f3: banana\n \n", patch)
}
