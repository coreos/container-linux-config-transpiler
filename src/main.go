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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/coreos/fuze/third_party/github.com/coreos/ignition/config"
	"github.com/coreos/fuze/third_party/github.com/go-yaml/yaml"
)

func stderr(f string, a ...interface{}) {
	out := fmt.Sprintf(f, a...)
	fmt.Fprintln(os.Stderr, strings.TrimSuffix(out, "\n"))
}

// hasUnrecognizedKeys finds unrecognized keys and warns about them on stderr.
// returns false when no unrecognized keys were found, true otherwise.
func hasUnrecognizedKeys(inCfg interface{}, refType reflect.Type) (warnings bool) {
	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
	}
	switch inCfg.(type) {
	case map[interface{}]interface{}:
		ks := inCfg.(map[interface{}]interface{})
	keys:
		for key := range ks {
			for i := 0; i < refType.NumField(); i++ {
				sf := refType.Field(i)
				tv := sf.Tag.Get("yaml")
				if tv == key {
					if warn := hasUnrecognizedKeys(ks[key], sf.Type); warn {
						warnings = true
					}
					continue keys
				}
			}

			stderr("Unrecognized keyword: %v", key)
			warnings = true
		}
	case []interface{}:
		ks := inCfg.([]interface{})
		for i := range ks {
			if warn := hasUnrecognizedKeys(ks[i], refType.Elem()); warn {
				warnings = true
			}
		}
	default:
	}
	return
}

func main() {
	flags := struct {
		help    bool
		pretty  bool
		inFile  string
		outFile string
	}{}

	flag.BoolVar(&flags.help, "help", false, "print help and exit")
	flag.BoolVar(&flags.pretty, "pretty", false, "indent the output file")
	flag.StringVar(&flags.inFile, "in-file", "/dev/stdin", "input file (YAML)")
	flag.StringVar(&flags.outFile, "out-file", "/dev/stdout", "output file (JSON)")

	flag.Parse()

	if flags.help {
		flag.Usage()
		return
	}

	cfg := config.Config{}
	dataIn, err := ioutil.ReadFile(flags.inFile)
	if err != nil {
		stderr("Failed to read: %v", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(dataIn, &cfg); err != nil {
		stderr("Failed to unmarshal input: %v", err)
		os.Exit(1)
	}

	var inCfg interface{}
	if err := yaml.Unmarshal(dataIn, &inCfg); err != nil {
		stderr("Failed to unmarshal input: %v", err)
		os.Exit(1)
	}

	if hasUnrecognizedKeys(inCfg, reflect.TypeOf(cfg)) {
		stderr("Unrecognized keys in input, aborting.")
		os.Exit(1)
	}

	var dataOut []byte
	if flags.pretty {
		dataOut, err = json.MarshalIndent(&cfg, "", "  ")
		dataOut = append(dataOut, '\n')
	} else {
		dataOut, err = json.Marshal(&cfg)
	}
	if err != nil {
		stderr("Failed to marshal output: %v", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(flags.outFile, dataOut, 0640); err != nil {
		stderr("Failed to write: %v", err)
		os.Exit(1)
	}
}
