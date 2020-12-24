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

package remote

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"goarrg.com/debug"
	"golang.org/x/crypto/openpgp"
)

func Get(url string, out io.WriteSeeker) error {
	resp, err := http.Get(url)

	if err != nil {
		return debug.ErrorWrap(err, "Failed to get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return debug.ErrorWrap(debug.ErrorNew("Bad status: %s", resp.Status), "Failed to get %s", url)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return debug.ErrorWrap(err, "Failed to get %s", url)
	}

	if _, err = out.Seek(0, io.SeekStart); err != nil {
		return debug.ErrorWrap(err, "Failed to get %s", url)
	}

	return debug.ErrorWrap(err, "Failed to get %s", url)
}

func verifyError(err error) error {
	return debug.ErrorWrap(err, "Verification failed")
}

func VerifyPGP(target, pk, sig io.ReadSeeker) error {
	keyring, err := openpgp.ReadArmoredKeyRing(pk)

	if err != nil {
		if _, err = pk.Seek(0, io.SeekStart); err != nil {
			return verifyError(err)
		}

		if keyring, err = openpgp.ReadKeyRing(pk); err != nil {
			return verifyError(err)
		}
	}

	if _, err = openpgp.CheckArmoredDetachedSignature(keyring, target, sig); err != nil {
		if _, err = target.Seek(0, io.SeekStart); err != nil {
			return verifyError(err)
		}
		if _, err = sig.Seek(0, io.SeekStart); err != nil {
			return verifyError(err)
		}

		if _, err = openpgp.CheckDetachedSignature(keyring, target, sig); err != nil {
			return verifyError(err)
		}
	}

	_, err = target.Seek(0, io.SeekStart)
	return verifyError(err)
}

func VerifySHA256(target io.ReadSeeker, targetSHA256 string) error {
	h := sha256.New()
	if _, err := io.Copy(h, target); err != nil {
		return verifyError(err)
	}

	sha256 := hex.EncodeToString(h.Sum(nil))
	if sha256 != targetSHA256 {
		return verifyError(debug.ErrorNew("SHA256 mismatch Got %q Want %q", sha256, targetSHA256))
	}

	_, err := target.Seek(0, io.SeekStart)
	return verifyError(err)
}
