// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package get

import (
	"encoding/json"
	"reflect"

	"fmt"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (TaskConfig) elemValueOrNil(v interface{}) interface{} {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Zero(t.Elem()).Interface()
		} else {
			return reflect.ValueOf(v).Interface()
		}
	} else if v == nil {
		return reflect.Zero(t).Interface()
	}

	return v
}

func (TaskConfig) mustMarshalJSON(v json.Marshaler) string {
	raw, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(raw)
}

// GetPFlagSet will return strongly types pflags for all fields in TaskConfig and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg TaskConfig) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("TaskResourceAttrConfig", pflag.ExitOnError)
	cmdFlags.StringVar(&(taskConfig.ExecFile), fmt.Sprintf("%v%v", prefix, "execFile"), taskConfig.ExecFile, "execution file name to be used for generating execution spec of a single task.")
	cmdFlags.StringVar(&(taskConfig.Version), fmt.Sprintf("%v%v", prefix, "version"), taskConfig.Version, "version of the task to be fetched.")
	cmdFlags.BoolVar(&(taskConfig.Latest), fmt.Sprintf("%v%v", prefix, "latest"), taskConfig.Latest, "flag to indicate to fetch the latest version, version flag will be ignored in this case")
	return cmdFlags
}
