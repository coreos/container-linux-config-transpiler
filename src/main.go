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
	"strings"

	"github.com/coreos/fuze/third_party/github.com/coreos/ignition/config"
	"github.com/coreos/fuze/third_party/github.com/go-yaml/yaml"
)

var (
	flagHelp    = flag.Bool("help", false, "print help and exit")
	flagInFile  = flag.String("in-file", "/dev/stdin", "input file (YAML)")
	flagOutFile = flag.String("out-file", "/dev/stdout", "output file (JSON)")
)

func stderr(f string, a ...interface{}) {
	out := fmt.Sprintf(f, a...)
	fmt.Fprintln(os.Stderr, strings.TrimSuffix(out, "\n"))
}

func stdout(f string, a ...interface{}) {
	out := fmt.Sprintf(f, a...)
	fmt.Fprintln(os.Stdout, strings.TrimSuffix(out, "\n"))
}

func panicf(f string, a ...interface{}) {
	panic(fmt.Sprintf(f, a...))
}

func main() {
	flag.Parse()

	if *flagHelp {
		flag.Usage()
		os.Exit(1)
	}

	cfg := config.Config{}
	b, err := ioutil.ReadFile(*flagInFile)
	if err != nil {
		stderr("Failed to read: %v", err)
		os.Exit(2)
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		stderr("Failed to unmarshal input: %v", err)
		os.Exit(3)
	}

	b, err = json.Marshal(&cfg)
	if err != nil {
		stderr("Failed to marshal output: %v", err)
		os.Exit(4)
	}

	if err := ioutil.WriteFile(*flagOutFile, b, 0640); err != nil {
		stderr("Failed to write: %v", err)
		os.Exit(5)
	}

	os.Exit(0)
}
