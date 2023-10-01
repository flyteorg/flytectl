package update

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v3"
)

// diffAsYaml marshals both objects as YAML and returns differences
// between marshalled values as a patch. Marshalling respects JSON
// field annotations.
func diffAsYaml(object1, object2 any) (string, error) {
	yaml1, err := marshalToYamlString(object1)
	if err != nil {
		return "", fmt.Errorf("diff as yaml: %w", err)
	}

	yaml2, err := marshalToYamlString(object2)
	if err != nil {
		return "", fmt.Errorf("diff as yaml: %w", err)
	}

	// TODO: kamal - remove
	fmt.Println("---")
	fmt.Println(yaml1)
	fmt.Println("---")
	fmt.Println(yaml2)
	fmt.Println("---")

	patch := diffStrings(yaml1, yaml2)
	return patch, nil
}

// marshalToYamlString marshals value to a YAML string, while respecting
// JSON field annotations.
func marshalToYamlString(value any) (string, error) {
	jsonText, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshalling object to json: %w", err)
	}

	var jsonObject interface{}
	if err := yaml.Unmarshal(jsonText, &jsonObject); err != nil {
		return "", fmt.Errorf("unmarshalling yaml to object: %w", err)
	}

	data, err := yaml.Marshal(jsonObject)
	if err != nil {
		return "", fmt.Errorf("marshalling object to yaml: %w", err)
	}

	return string(data), nil
}

// diffStrings returns differences between two strings as a patch.
func diffStrings(s1, s2 string) string {
	dmp := diffmatchpatch.New()

	lines1, lines2, lineIndex := dmp.DiffLinesToRunes(s1, s2)

	// checklines: Speedup flag. If false, then don't run a line-level
	// diff first to identify the changed areas. If true, then run a
	// faster slightly less optimal diff.
	// (Remove this comment once the documentation has this integrated:
	//  https://github.com/sergi/go-diff/issues/95)
	diffs := dmp.DiffMainRunes(lines1, lines2, false /* checklines */)
	equal := len(diffs) == 0 ||
		len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual
	if equal {
		return ""
	}

	diffs = dmp.DiffCharsToLines(diffs, lineIndex)

	patches := dmp.PatchMake(s1, diffs)
	patchesText := dmp.PatchToText(patches)

	// There is weird behavior (seemingly a bug) which nobody knows
	// the reason for: original DMP implementation URL-escapes characters,
	// and then unescapes only some of them. LF is one of the characters
	// that gets left escaped as '%0A', yet unescaped LF also ends up in
	// the output string. Thus, we strip the '%0A'.
	//
	// https://github.com/sergi/go-diff/issues/87
	// https://github.com/google/diff-match-patch/issues/4
	patchesText = strings.ReplaceAll(patchesText, "%0A", "\n")

	// TODO: kamal - there are other unescaped characters

	return patchesText
}
