// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package project

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var dereferencableKindsProjectCreateConfig = map[reflect.Kind]struct{}{
	reflect.Array: {}, reflect.Chan: {}, reflect.Map: {}, reflect.Ptr: {}, reflect.Slice: {},
}

// Checks if t is a kind that can be dereferenced to get its underlying type.
func canGetElementProjectCreateConfig(t reflect.Kind) bool {
	_, exists := dereferencableKindsProjectCreateConfig[t]
	return exists
}

// This decoder hook tests types for json unmarshaling capability. If implemented, it uses json unmarshal to build the
// object. Otherwise, it'll just pass on the original data.
func jsonUnmarshalerHookProjectCreateConfig(_, to reflect.Type, data interface{}) (interface{}, error) {
	unmarshalerType := reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	if to.Implements(unmarshalerType) || reflect.PtrTo(to).Implements(unmarshalerType) ||
		(canGetElementProjectCreateConfig(to.Kind()) && to.Elem().Implements(unmarshalerType)) {

		raw, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Failed to marshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		res := reflect.New(to).Interface()
		err = json.Unmarshal(raw, &res)
		if err != nil {
			fmt.Printf("Failed to umarshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		return res, nil
	}

	return data, nil
}

func decode_ProjectCreateConfig(input, result interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		Result:           result,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			jsonUnmarshalerHookProjectCreateConfig,
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func join_ProjectCreateConfig(arr interface{}, sep string) string {
	listValue := reflect.ValueOf(arr)
	strs := make([]string, 0, listValue.Len())
	for i := 0; i < listValue.Len(); i++ {
		strs = append(strs, fmt.Sprintf("%v", listValue.Index(i)))
	}

	return strings.Join(strs, sep)
}

func testDecodeJson_ProjectCreateConfig(t *testing.T, val, result interface{}) {
	assert.NoError(t, decode_ProjectCreateConfig(val, result))
}

func testDecodeRaw_ProjectCreateConfig(t *testing.T, vStringSlice, result interface{}) {
	assert.NoError(t, decode_ProjectCreateConfig(vStringSlice, result))
}

func TestProjectCreateConfig_GetPFlagSet(t *testing.T) {
	val := ProjectCreateConfig{}
	cmdFlags := val.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())
}

func TestProjectCreateConfig_SetFlags(t *testing.T) {
	actual := ProjectCreateConfig{}
	cmdFlags := actual.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())

	t.Run("Test_id", func(t *testing.T) {

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("id", testValue)
			if vString, err := cmdFlags.GetString("id"); err == nil {
				testDecodeJson_ProjectCreateConfig(t, fmt.Sprintf("%v", vString), &actual.ID)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_name", func(t *testing.T) {

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("name", testValue)
			if vString, err := cmdFlags.GetString("name"); err == nil {
				testDecodeJson_ProjectCreateConfig(t, fmt.Sprintf("%v", vString), &actual.Name)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_file", func(t *testing.T) {

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("file", testValue)
			if vString, err := cmdFlags.GetString("file"); err == nil {
				testDecodeJson_ProjectCreateConfig(t, fmt.Sprintf("%v", vString), &actual.File)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_description", func(t *testing.T) {

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("description", testValue)
			if vString, err := cmdFlags.GetString("description"); err == nil {
				testDecodeJson_ProjectCreateConfig(t, fmt.Sprintf("%v", vString), &actual.Description)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_labels", func(t *testing.T) {

		t.Run("Override", func(t *testing.T) {
			testValue := "a=1,b=2"

			cmdFlags.Set("labels", testValue)
			if vStringToString, err := cmdFlags.GetStringToString("labels"); err == nil {
				testDecodeRaw_ProjectCreateConfig(t, vStringToString, &actual.Labels)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_dryRun", func(t *testing.T) {

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("dryRun", testValue)
			if vBool, err := cmdFlags.GetBool("dryRun"); err == nil {
				testDecodeJson_ProjectCreateConfig(t, fmt.Sprintf("%v", vBool), &actual.DryRun)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
}
