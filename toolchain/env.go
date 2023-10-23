/*
Copyright 2022 The goARRG Authors.

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
	"fmt"
	"os"
	"runtime"
	"sort"
)

var envMap = map[string]struct{}{}

/*
EnvRegister registers the env var for printing and sets it to value if it has
no current value.
*/
func EnvRegister(key, value string) {
	currentValue := os.Getenv(key)
	if currentValue == "" && value != "" {
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
	envMap[key] = struct{}{}
}

/*
EnvGet registers the env var for printing and returns the current value.
*/
func EnvGet(key string) string {
	envMap[key] = struct{}{}
	return os.Getenv(key)
}

/*
EnvSet registers the env var for printing and always sets it to value.
*/
func EnvSet(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		panic(err)
	}
	envMap[key] = struct{}{}
}

/*
EnvString returns a string containing every registered env var and its' value.
*/
func EnvString() string {
	envOut := ""
	envList := make([]string, 0, len(envMap))
	for k := range envMap {
		envList = append(envList, k)
	}
	sort.Strings(envList)
	for _, k := range envList {
		if runtime.GOOS == "windows" {
			envOut += fmt.Sprintf("set %s=%s\n", k, os.Getenv(k))
		} else {
			envOut += fmt.Sprintf("%s=%s\n", k, os.Getenv(k))
		}
	}
	return envOut[:len(envOut)-1]
}
