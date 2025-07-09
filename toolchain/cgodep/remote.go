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

package cgodep

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"goarrg.com/debug"
)

/*
Get downloads the file located at url and places it within filepath.Join(DataDir(), cache)
and names it fileName, verify will be used to verify the contents. If filename is
found within the cache, verify is used to verify the contents and only redownloading if errored.
*/
func Get(url string, fileName string, verify func(io.ReadSeeker) error) (io.ReadSeekCloser, error) {
	cacheDir := CacheDir()
	fileName = filepath.Join(cacheDir, fileName)

	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, debug.ErrorWrapf(err, "Failed to create download cache dir: %q", cacheDir)
	}

	file, err := os.OpenFile(fileName, os.O_RDWR, 0)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, debug.ErrorWrapf(err, "Failed to read cached file: %q", fileName)
		}
		file, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0o655)
		if err != nil {
			return nil, err
		}
	} else {
		debug.VPrintf("Verifying cached file: %q", fileName)
		if err := verify(file); err == nil {
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				file.Close()
				return nil, err
			}
			return file, nil
		}
		if err := file.Truncate(0); err != nil {
			file.Close()
			return nil, err
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			file.Close()
			return nil, err
		}
	}

	debug.VPrintf("Downloading: %q", url)

	resp, err := http.Get(url)
	if err != nil {
		file.Close()
		return nil, debug.ErrorWrapf(err, "Failed to get: %q", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		file.Close()
		return nil, debug.ErrorWrapf(debug.Errorf("Bad status: %s", resp.Status), "Failed to get %s", url)
	}
	if _, err := io.Copy(file, resp.Body); err != nil {
		file.Close()
		return nil, debug.ErrorWrapf(err, "Failed to get: %q", url)
	}

	debug.VPrintf("Verifying downloaded data")

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		file.Close()
		return nil, err
	}
	if err := verify(file); err != nil {
		file.Close()
		return nil, debug.ErrorWrapf(err, "Failed to verify downloaded data")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		file.Close()
		return nil, err
	}
	return file, nil
}

/*
VerifySHA256 is a convenience function to verify a sha256 checksum.
*/
func VerifySHA256(target io.ReadSeeker, targetSHA256 string) error {
	h := sha256.New()
	if _, err := io.Copy(h, target); err != nil {
		return err
	}
	sha256 := hex.EncodeToString(h.Sum(nil))
	if sha256 != targetSHA256 {
		return debug.Errorf("SHA256 mismatch, Got: %q Want: %q", sha256, targetSHA256)
	}
	return nil
}
