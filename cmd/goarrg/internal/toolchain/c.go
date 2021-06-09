/*
Copyright 2021 The goARRG Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package toolchain

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"goarrg.com/debug"
)

func FindHeader(cc, header string) (string, error) {
	ex := exec.Command(cc, "-M", "-E", "-")
	ex.Env = os.Environ()
	ex.Stdin = bytes.NewReader([]byte("#include<" + header + ">"))
	out, err := ex.CombinedOutput()
	if err != nil {
		return "", debug.ErrorNew("Failed to find %q using %q: %q", header, cc, string(out))
	}

	for _, s := range strings.Fields(string(out)) {
		if strings.Contains(s, header) {
			return filepath.Dir(s), nil
		}
	}

	// should never be here
	panic(debug.ErrorNew("%q found %q but unable to find header in output: %q", cc, header, string(out)))
}
