/*
Copyright 2026 The goARRG Authors.

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

package web

import (
	"encoding/json"
	"io"
	"net/http"

	"goarrg.com/debug"
)

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, debug.ErrorWrapf(err, "Failed to get: %q", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, debug.ErrorWrapf(debug.Errorf("Bad status: %s", resp.Status), "Failed to get %s", url)
	}

	return io.ReadAll(resp.Body)
}

func GetJSON(url string, out any) error {
	data, err := Get(url)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, out)
}
