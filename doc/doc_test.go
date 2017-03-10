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

package doc

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/coreos/container-linux-config-transpiler/config"
)

// TestValidYAML reads in all .md files in the current directory, finds parts of
// the file wrapped in ```yaml and ```, and attempts to parse/validate it with
// ct. If the report is >0 this test fails (it fails on warnings).
func TestValidYAML(t *testing.T) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		t.Errorf("couldn't read dir: %v", err)
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".md") {
			// Only check markdown files
			continue
		}
		fileContents, err := ioutil.ReadFile(f.Name())
		if err != nil {
			t.Errorf("couldn't read file %s: %v", f.Name(), err)
		}

		fileLines := strings.Split(string(fileContents), "\n")

		yamlSections := findYamlSections(fileLines)

		for _, yaml := range yamlSections {
			_, report := config.Parse([]byte(strings.Join(yaml, "\n")))
			reportStr := report.String()
			if reportStr != "" {
				t.Errorf("non-empty parsing report in %s: %s", f.Name(), reportStr)
			}
		}
	}
}

func findYamlSections(fileLines []string) [][]string {
	var yamlSections [][]string
	var currentSection []string
	inASection := false
	for _, line := range fileLines {
		if line == "```" {
			inASection = false
			if len(currentSection) > 0 {
				yamlSections = append(yamlSections, currentSection)
			}
			currentSection = nil
		}
		if inASection {
			currentSection = append(currentSection, line)
		}
		if line == "```yaml" {
			inASection = true
		}
	}
	return yamlSections
}
