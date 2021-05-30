// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package workflow

import (
	"encoding/json"
	"reflect"

	"fmt"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (Config) elemValueOrNil(v interface{}) interface{} {
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

func (Config) mustJsonMarshal(v interface{}) string {
	raw, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(raw)
}

func (Config) mustMarshalJSON(v json.Marshaler) string {
	raw, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(raw)
}

// GetPFlagSet will return strongly types pflags for all fields in Config and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg Config) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("Config", pflag.ExitOnError)
	cmdFlags.StringVar(&DefaultConfig.Version, fmt.Sprintf("%v%v", prefix, "version"), DefaultConfig.Version, "version of the workflow to be fetched.")
	cmdFlags.BoolVar(&DefaultConfig.Latest, fmt.Sprintf("%v%v", prefix, "latest"), DefaultConfig.Latest, " flag to indicate to fetch the latest version,  version flag will be ignored in this case")
	cmdFlags.StringVar(&DefaultConfig.OutputFileName, fmt.Sprintf("%v%v", prefix, "opFileName"), DefaultConfig.OutputFileName, "output file name to be used for the fetched workflow.Currently support dumping the SVG image.Use .svg extension for the filename. Defaults to /tmp/workflow.svg")
	return cmdFlags
}
