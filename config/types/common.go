// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"fmt"
	"reflect"
)

func isZero(v interface{}) bool {
	if v == nil {
		return true
	}
	zv := reflect.Zero(reflect.TypeOf(v))
	return reflect.DeepEqual(v, zv.Interface())
}

// serviceContentsFromEnvVars builds the systemd drop in from a list of ENV_VAR=VALUE strings.
func serviceContentsFromEnvVars(vars []string) string {
	out := "[Service]\n"
	for _, v := range vars {
		out += fmt.Sprintf("Environment=\"%s\"\n", v)
	}
	return out
}

// getEnvVars builds a list of ENV_VAR=VALUE from a struct with env: tags on its members.
func getEnvVars(e interface{}) []string {
	if e == nil {
		return nil
	}
	et := reflect.TypeOf(e)
	ev := reflect.ValueOf(e)

	vars := []string{}
	for i := 0; i < et.NumField(); i++ {
		if val := ev.Field(i).Interface(); !isZero(val) {
			if et.Field(i).Anonymous {
				vars = append(vars, getEnvVars(val)...)
			} else {
				key := et.Field(i).Tag.Get("env")
				vars = append(vars, fmt.Sprintf("%s=%v", key, val))
			}
		}
	}

	return vars
}
