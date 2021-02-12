// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package register

import (
	"encoding/json"
	"reflect"

	"fmt"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (FilesConfig) elemValueOrNil(v interface{}) interface{} {
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

func (FilesConfig) mustMarshalJSON(v json.Marshaler) string {
	raw, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(raw)
}

// GetPFlagSet will return strongly types pflags for all fields in FilesConfig and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg FilesConfig) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("FilesConfig", pflag.ExitOnError)
	cmdFlags.StringVarP(&(filesConfig.Version),fmt.Sprintf("%v%v", prefix, "version"), "v", "v1", "version of the entity to be registered with flyte.")
	cmdFlags.BoolVarP(&(filesConfig.SkipOnError), fmt.Sprintf("%v%v", prefix, "skipOnError"), "s", *new(bool), "fail fast when registering files.")
	cmdFlags.BoolVarP(&(filesConfig.Archive), fmt.Sprintf("%v%v", prefix, "archive"), "a", *new(bool), "pass in archive file either an http link or local path.")
	return cmdFlags
}
