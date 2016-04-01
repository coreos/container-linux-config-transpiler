// Copyright 2015 CoreOS, Inc.
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

package main

import (
	"reflect"
	"testing"

	"github.com/coreos/ignition/config"
	"github.com/go-yaml/yaml"
)

func TestHasUnrecognizedKeys(t *testing.T) {
	tests := []struct {
		in    string
		unrec bool
	}{
		{
			in:    "ignition_version: 1",
			unrec: false,
		},
		{
			in:    "ignition_version: 1\npasswd:\n users:\n  - name: foobar\n",
			unrec: false,
		},
		{
			in:    "foo: bar",
			unrec: true,
		},
		{
			in:    "ignition_version: 1\npasswd:\n users:\n  - naem: foobar\n",
			unrec: true,
		},
	}

	for i, tt := range tests {
		var cfg interface{}
		if err := yaml.Unmarshal([]byte(tt.in), &cfg); err != nil {
			t.Errorf("%d: unmarshal failed: %v", i, err)
			continue
		}
		if unrec := hasUnrecognizedKeys(cfg, reflect.TypeOf(config.Config{})); unrec != tt.unrec {
			t.Errorf("%d: expected %v got %v", i, tt.unrec, unrec)
		}
	}
}
