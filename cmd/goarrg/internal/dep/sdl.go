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

package dep

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/cmd/goarrg/internal/remote"
)

const sdlVR = "SDL2-2.0.14"

const sdlPK = `
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1.0.4 (GNU/Linux)
Comment: For info see http://www.gnupg.org

mQGiBDpWOb0RBADQwd3d9mzt6KzqlsgXf9mikBuMbpKzYs1SBKYpdzUs9sRY0CnH
vCQTrL5sI57yKLnqEl6SbIiE75ZwrSWwvUDFFTh35Jew5nPZwv64en2kw2y4qrnJ
kBZCHDSU4KgfUZtoJ25Tmeru5MLNbXxCOoMszO5L5OchwMrGMtmFLRA/bwCgy5Th
d1/vJo+bej9tbgv++SJ05o0D/3MPK7EBoxWkQ0I+ScqOsvSMRQXWc/hXy4lyIp8e
xJByBApkv0LiiT3KlPpq/K2gTlDlCZ/JTt6Rv8Ug0g47R3a0aoz9kfc15UjHdiap
UOfF9MWmmbw59Lyx6+y2e0/C5xWzNOR1G4G5y4RZL/GXrp67xz/0fEhI85R+eASq
AEfSBAC5ZxwnBwyl+h+PXeJYKrPQjSUlgtSAkKp7PNBywwlue1LcSb7j4cc+cmgH
QMVuM883LPE59btNzFTAZjlzzIMiaXf5h9EkDARTGQ1wFiO3V5vIbVLh4kAoNfpT
egy7bYn3UrlbKg3V2DbCdEXm1zQufZzK7T0yenA5Ps8xXX7mNrQhU2FtIExhbnRp
bmdhIDxzbG91a2VuQGxpYnNkbC5vcmc+iFcEExECABcFAjpWOb0FCwcKAwQDFQMC
AxYCAQIXgAAKCRAwpZN3p3Y75t9RAJ48WI+nOPes0WK7t381Ij4JfSYxWQCgjpMa
Dg3/ah23HZhYtTKtHUzD9zi5AQ0EOlY5wxAEAPvjB0B5RNAj8hBF/Lq78w5rJ1/f
5RqWXmdfxApuEE/9OEFXUSUXms9f/IWvySdyf48Pk4t2h8b8i7F0f3R+tcCp6m0P
t1BSNHYumfmtonTy5FHqpwBVlEi7I0s5mD3kxO+k8PQbATHH5smFnoz2UTc+MzQj
UdtTzXUkUgqvf9zTAAMGA/9Y/h6rhi3YYXeI6SmbXqcmzsQKzaWVhLew67szejnY
sKIJ1ja4MefYlthCXgmIBriNftxIGtBI0Pcmzwpn0eknRNK6NgpmESbGKCWh59Je
iAK5hdBPe47LSFVct5zSO9vQhRDyLzhzPPtB3XeoKTUkLWxBSLbeZVwcHPIK/wLa
l4hGBBgRAgAGBQI6VjnDAAoJEDClk3endjvmxmUAn3Ne6Z3UULpie8RJP15RBt6K
2MWFAJ9hK/Ls/FeBJ9d50qxmYdZ2RrTXNg==
=toqC
-----END PGP PUBLIC KEY BLOCK-----
`

var sdlSig = []byte{
	0x88, 0x5d, 0x04, 0x00, 0x11, 0x02, 0x00, 0x1d, 0x16, 0x21, 0x04, 0x15,
	0x28, 0x63, 0x5d, 0x80, 0x53, 0xa5, 0x7f, 0x77, 0xd1, 0xe0, 0x86, 0x30,
	0xa5, 0x93, 0x77, 0xa7, 0x76, 0x3b, 0xe6, 0x05, 0x02, 0x5f, 0xe0, 0xe0,
	0x72, 0x00, 0x0a, 0x09, 0x10, 0x30, 0xa5, 0x93, 0x77, 0xa7, 0x76, 0x3b,
	0xe6, 0xfe, 0x25, 0x00, 0x9b, 0x05, 0x38, 0x4b, 0xdf, 0xbd, 0x7d, 0x16,
	0x30, 0x51, 0x9c, 0xfc, 0x2a, 0xf1, 0xb7, 0x41, 0xf9, 0xd4, 0x53, 0xda,
	0x20, 0x00, 0x9e, 0x2f, 0xf7, 0x9d, 0x49, 0x70, 0xf2, 0x33, 0x84, 0x92,
	0x27, 0xd2, 0x8f, 0x0d, 0xff, 0x40, 0x69, 0xde, 0xd7, 0xdd, 0xb6,
}

func init() {
	var installFn, prebuildFn func()
	var config, staticConfig string

	switch base.GOOS() {
	case "linux":
		installFn = sdlLinux
		config = sdlConfigLinux
		staticConfig = sdlStaticConfigLinux
	case "windows":
		installFn = sdlWindows
		config = sdlConfigWindows
		staticConfig = sdlStaticConfigWindows
	default:
		return
	}

	deps["sdl"] = dep{
		sdlVR + ".tar.gz", "https://www.libsdl.org/release/" + sdlVR + ".tar.gz", sdlVR + "-0",
		func(r io.ReadSeeker) error {
			return remote.VerifyPGP(r, bytes.NewReader([]byte(sdlPK)), bytes.NewReader(sdlSig))
		},
		func() {
			if err := os.MkdirAll("build", 0o755); err != nil {
				panic(err)
			}

			if err := os.Chdir("build"); err != nil {
				panic(err)
			}

			installFn()

			sdlPath := filepath.ToSlash(usrData)
			err := ioutil.WriteFile(usrData+"/lib/pkgconfig/sdl2.pc", []byte("prefix="+sdlPath+
				"\nexec_prefix=${prefix}"+
				"\nlibdir=${exec_prefix}/lib"+
				"\nincludedir=${prefix}/include"+
				"\n"+
				"\nName: sdl2"+
				"\nDescription:"+
				"\nVersion: "+sdlVR+
				"\nRequires:"+
				"\nConflicts:"+
				"\n"+config+
				"\n"), 0o644)

			if err != nil {
				panic(err)
			}

			err = ioutil.WriteFile(usrData+"/lib/pkgconfig/sdl2-static.pc", []byte("prefix="+sdlPath+
				"\nexec_prefix=${prefix}"+
				"\nlibdir=${exec_prefix}/lib"+
				"\nincludedir=${prefix}/include"+
				"\n"+
				"\nName: sdl2-static"+
				"\nDescription:"+
				"\nVersion: "+sdlVR+
				"\nRequires:"+
				"\nConflicts: sdl2"+
				"\n"+staticConfig+
				"\n"), 0o644)

			if err != nil {
				panic(err)
			}
		},
		prebuildFn,
	}
}
