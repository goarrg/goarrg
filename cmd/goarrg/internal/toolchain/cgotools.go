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
	"os"
	"path/filepath"
	"runtime"

	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/exec"
)

func setupCgoTools() {
	toolsDir := filepath.Join(cmd.ModuleDataPath(), "tools")
	tool := "goarrg-config"
	toolFile := filepath.Join(toolsDir, tool)

	if runtime.GOOS == "windows" {
		toolFile += ".exe"
	}

	if err := exec.Run("go", "build", "-o", toolFile, "-ldflags=-X main.dataDir="+cmd.ModuleDataPath(), "goarrg.com/cmd/tool/"+tool); err != nil {
		panic(err)
	}

	if err := os.Setenv("PKG_CONFIG", toolFile); err != nil {
		panic(err)
	}
}
