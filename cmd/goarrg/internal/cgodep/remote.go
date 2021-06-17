/*
Copyright 2020 The goARRG Authors.

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
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"goarrg.com/debug"
	"golang.org/x/crypto/openpgp"
)

func get(url string, verify func([]byte) error) ([]byte, error) {
	cacheDir := cgoDepCache()
	fileName := filepath.Join(cacheDir, filepath.Base(url))

	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		panic(debug.ErrorWrapf(err, "Failed to create download cache dir: %q", cacheDir))
	}

	if data, err := os.ReadFile(fileName); err != nil {
		if !os.IsNotExist(err) {
			return nil, debug.ErrorWrapf(err, "Failed to read cached file: %q", fileName)
		}
	} else {
		debug.LogV("Verifying cached file: %q", fileName)
		if err := verify(data); err == nil {
			return data, nil
		}
	}

	debug.LogV("Downloading: %q", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, debug.ErrorWrapf(err, "Failed to get: %q", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, debug.ErrorWrapf(debug.Errorf("Bad status: %s", resp.Status), "Failed to get %s", url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, debug.ErrorWrapf(err, "Failed to get: %q", url)
	}

	debug.LogV("Verifying downloaded data")

	if err := verify(data); err != nil {
		return nil, debug.ErrorWrapf(err, "Failed to verify downloaded data")
	}

	debug.LogV("Writing to: %q", fileName)

	if err := os.WriteFile(fileName, data, 0o644); err != nil {
		return nil, debug.ErrorWrapf(err, "Failed to write file: %q", fileName)
	}

	return data, nil
}

func verifyError(err error) error {
	return debug.ErrorWrapf(err, "Verification failed")
}

func verifyPGP(target, pk, sig []byte) error {
	keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(pk))
	if err != nil {
		if keyring, err = openpgp.ReadKeyRing(bytes.NewReader(pk)); err != nil {
			return verifyError(err)
		}
	}

	if _, err = openpgp.CheckArmoredDetachedSignature(keyring, bytes.NewReader(target), bytes.NewReader(sig)); err != nil {
		if _, err = openpgp.CheckDetachedSignature(keyring, bytes.NewReader(target), bytes.NewReader(sig)); err != nil {
			return verifyError(err)
		}
	}

	return nil
}

func verifySHA256(target []byte, targetSHA256 string) error {
	h := sha256.New()
	if _, err := io.Copy(h, bytes.NewReader(target)); err != nil {
		return verifyError(err)
	}

	sha256 := hex.EncodeToString(h.Sum(nil))
	if sha256 != targetSHA256 {
		return verifyError(debug.Errorf("SHA256 mismatch, Got: %q Want: %q", sha256, targetSHA256))
	}

	return nil
}
