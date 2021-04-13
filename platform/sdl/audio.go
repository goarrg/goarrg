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

package sdl

/*
	#cgo pkg-config: sdl2
	#include <SDL2/SDL.h>

	static SDL_RWops* RWFromUintptr(uintptr_t ptr, int sz) {
		return SDL_RWFromConstMem((void*)ptr, sz);
	}
*/
import "C"
import (
	"reflect"
	"unsafe"

	"goarrg.com"
	"goarrg.com/asset"
	"goarrg.com/asset/audio"
	"goarrg.com/debug"
)

/*
	The default SDL channel interleaving order is:
	1: L                       (mono)
	2: L R                     (stereo)
	3: L R LF                  (2.1 surround)
	4: L R BSL BSR             (quad)
	5: L R C BSL BSR           (quad + center)
	6: L R C LF SL SR          (5.1 surround)
	7: L R C LF BC SL SR       (6.1 surround)
	8: L R C LF BSL BSR SL SR  (7.1 surround)
*/

type audioChannelMask audio.Channel

const (
	AudioChannelBackCenter audio.Channel = audio.ChannelCount
	AudioChannelCount                    = audio.ChannelCount + 1

	audioChannelMaskMono            = audioChannelMask(1 << audio.ChannelLeft)
	audioChannelMaskStereo          = audioChannelMaskMono | audioChannelMask(1<<audio.ChannelRight)
	audioChannelMask2Point1         = audioChannelMaskStereo | audioChannelMask(1<<audio.ChannelLowFrequency)
	audioChannelMaskQuad            = audioChannelMaskStereo | audioChannelMask((1<<audio.ChannelBackSurroundLeft)|(1<<audio.ChannelBackSurroundRight))
	audioChannelMaskQuadCenter      = audioChannelMaskQuad | audioChannelMask(1<<audio.ChannelCenter)
	audioChannelMaskSurround5Point1 = audioChannelMask2Point1 | audioChannelMask((1<<audio.ChannelCenter)|(1<<audio.ChannelSurroundLeft)|(1<<audio.ChannelSurroundRight))
	audioChannelMask6Point1         = audioChannelMaskSurround5Point1 | audioChannelMask(1<<AudioChannelBackCenter)
	audioChannelMaskSurround7Point1 = audioChannelMaskSurround5Point1 | audioChannelMask((1<<audio.ChannelBackSurroundLeft)|(1<<audio.ChannelBackSurroundRight))
)

type audioSystem struct {
	cfg   goarrg.AudioConfig
	mixer goarrg.Audio
	cAID  C.SDL_AudioDeviceID
	buf   []float32
}

func channelCountToList(n int) ([]audio.Channel, error) {
	switch n {
	case 1:
		return []audio.Channel{audio.ChannelLeft}, nil
	case 2:
		return []audio.Channel{audio.ChannelLeft, audio.ChannelRight}, nil
	case 3:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelLowFrequency,
		}, nil
	case 4:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelBackSurroundLeft,
			audio.ChannelBackSurroundRight,
		}, nil
	case 5:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelCenter,
			audio.ChannelBackSurroundLeft,
			audio.ChannelBackSurroundRight,
		}, nil
	case 6:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelCenter,
			audio.ChannelLowFrequency,
			audio.ChannelSurroundLeft,
			audio.ChannelSurroundRight,
		}, nil
	case 7:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelCenter,
			audio.ChannelLowFrequency,
			AudioChannelBackCenter,
			audio.ChannelSurroundLeft,
			audio.ChannelSurroundRight,
		}, nil
	case 8:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelCenter,
			audio.ChannelLowFrequency,
			audio.ChannelBackSurroundLeft,
			audio.ChannelBackSurroundRight,
			audio.ChannelSurroundLeft,
			audio.ChannelSurroundRight,
		}, nil
	default:
		return nil, debug.ErrorNew("Unsupported channel value %d", n)
	}
}

func verifyChannelList(list []audio.Channel) ([]audio.Channel, error) {
	m := audioChannelMask(0)

	for _, c := range list {
		m = m | audioChannelMask(1<<c)
	}

	switch m {
	case audioChannelMaskMono:
		return []audio.Channel{audio.ChannelLeft}, nil
	case audioChannelMaskStereo:
		return []audio.Channel{audio.ChannelLeft, audio.ChannelRight}, nil
	case audioChannelMask2Point1:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelLowFrequency,
		}, nil
	case audioChannelMaskQuad:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelBackSurroundLeft,
			audio.ChannelBackSurroundRight,
		}, nil
	case audioChannelMaskQuadCenter:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelCenter,
			audio.ChannelBackSurroundLeft,
			audio.ChannelBackSurroundRight,
		}, nil
	case audioChannelMaskSurround5Point1:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelCenter,
			audio.ChannelLowFrequency,
			audio.ChannelSurroundLeft,
			audio.ChannelSurroundRight,
		}, nil
	case audioChannelMask6Point1:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelCenter,
			audio.ChannelLowFrequency,
			AudioChannelBackCenter,
			audio.ChannelSurroundLeft,
			audio.ChannelSurroundRight,
		}, nil
	case audioChannelMaskSurround7Point1:
		return []audio.Channel{
			audio.ChannelLeft,
			audio.ChannelRight,
			audio.ChannelCenter,
			audio.ChannelLowFrequency,
			audio.ChannelBackSurroundLeft,
			audio.ChannelBackSurroundRight,
			audio.ChannelSurroundLeft,
			audio.ChannelSurroundRight,
		}, nil
	default:
		return nil, debug.ErrorNew("Unsupported channel list %v", list)
	}
}

func decodeWAV(a asset.Asset) (audio.Spec, int, []float32, error) {
	cR := C.RWFromUintptr(C.uintptr_t(a.Uintptr()), C.int(a.Size()))
	cSpec := C.SDL_AudioSpec{}
	cBuf := (*C.Uint8)(nil)
	cLen := C.Uint32(0)

	//nolint:staticcheck
	if C.SDL_LoadWAV_RW(cR, 0, &cSpec, &cBuf, &cLen) == nil {
		err := debug.ErrorWrap(debug.ErrorNew(C.GoString(C.SDL_GetError())), "Failed to decode WAV")
		C.SDL_ClearError()
		return audio.Spec{}, 0, nil, err
	}

	defer C.SDL_FreeWAV(cBuf)

	channels, err := channelCountToList(int(cSpec.channels))

	if err != nil {
		return audio.Spec{}, 0, nil, debug.ErrorWrap(err, "Failed to decode WAV")
	}

	var cCVT = C.SDL_AudioCVT{}

	//nolint:staticcheck
	if C.SDL_BuildAudioCVT(&cCVT, cSpec.format, cSpec.channels, cSpec.freq, C.AUDIO_F32SYS, cSpec.channels, cSpec.freq) < 0 {
		err := debug.ErrorWrap(debug.ErrorNew(C.GoString(C.SDL_GetError())), "Failed to decode WAV")
		C.SDL_ClearError()
		return audio.Spec{}, 0, nil, err
	}

	samples := int(cLen) / (int(cSpec.format&0xFF) / 8)
	cTrack := *(*[]float32)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(cBuf)), samples, samples,
	}))

	if cCVT.needed == 1 {
		sz := samples * int(unsafe.Sizeof(float32(0)))
		cCVT.len = C.int(cLen)
		cCVT.buf = (*C.Uint8)(C.malloc(C.size_t(cCVT.len * cCVT.len_mult)))

		defer C.free(unsafe.Pointer(cCVT.buf))

		C.memmove(unsafe.Pointer(cCVT.buf), unsafe.Pointer(cBuf), C.size_t(cLen))

		//nolint:staticcheck
		if C.SDL_ConvertAudio(&cCVT) < 0 {
			err := debug.ErrorWrap(debug.ErrorNew(C.GoString(C.SDL_GetError())), "Failed to decode WAV")
			C.SDL_ClearError()
			return audio.Spec{}, 0, nil, err
		}

		if sz != int(cCVT.len_cvt) {
			return audio.Spec{}, 0, nil, debug.ErrorWrap(
				debug.ErrorNew("Calculated buffer size does not match SDL size, Got: %d SDL: %d",
					sz, int(cCVT.len_cvt)), "Failed to decode WAV")
		}

		cTrack = *(*[]float32)(unsafe.Pointer(&reflect.SliceHeader{
			uintptr(unsafe.Pointer(cCVT.buf)), samples, samples,
		}))
	}

	return audio.Spec{
			Frequency: int(cSpec.freq),
			Channels:  channels,
		},
		samples,
		append([]float32(nil), cTrack...), nil
}

func (*platform) AudioInit(mixer goarrg.Audio) error {
	if Platform.config.Audio.Importer.EnableWAV {
		audio.RegisterFormat("RIFF????WAVE", decodeWAV)
	}

	if mixer == nil {
		debug.LogI("SDL audio disabled")
		return nil
	}

	if C.SDL_InitSubSystem(C.SDL_INIT_AUDIO) != 0 {
		err := debug.ErrorWrap(debug.ErrorNew(C.GoString(C.SDL_GetError())), "Failed to init SDL audio")
		C.SDL_ClearError()
		return err
	}

	cfg := mixer.AudioConfig()

	if len(cfg.Spec.Channels) == 0 {
		return debug.ErrorWrap(debug.ErrorNew("No channels defined"), "Failed to init SDL audio")
	}

	if cfg.BufferLength == 0 || (cfg.BufferLength&(cfg.BufferLength-1)) != 0 {
		return debug.ErrorWrap(debug.ErrorNew("BufferLength must be a power of 2"), "Failed to init SDL audio")
	}

	channels, err := verifyChannelList(cfg.Spec.Channels)

	if err != nil {
		return debug.ErrorWrap(err, "Failed to init SDL audio")
	}

	cfg.Spec.Channels = channels

	want := C.SDL_AudioSpec{
		freq:     C.int(cfg.Spec.Frequency),
		format:   C.AUDIO_F32SYS,
		channels: C.Uint8(len(channels)),
		samples:  C.Uint16(cfg.BufferLength),
	}
	got := C.SDL_AudioSpec{}

	//nolint:staticcheck
	cAID := C.SDL_OpenAudioDevice(nil, 0, &want, &got, 0)

	if cAID == 0 {
		err := debug.ErrorWrap(debug.ErrorNew(C.GoString(C.SDL_GetError())), "Failed to init SDL audio")
		C.SDL_ClearError()
		return err
	}

	cfg.Spec.Channels, err = channelCountToList(int(got.channels))

	if err != nil {
		return debug.ErrorWrap(err, "Failed to init SDL audio")
	}

	cfg.Spec.Frequency = int(got.freq)
	cfg.BufferLength = int(got.samples)

	if err := mixer.Init(cfg); err != nil {
		return debug.ErrorWrap(err, "Failed to init SDL audio")
	}

	Platform.audio.cAID = cAID
	Platform.audio.mixer = mixer
	Platform.audio.cfg = cfg
	Platform.audio.buf = make([]float32, cfg.Spec.Frequency*len(cfg.Spec.Channels))

	C.SDL_PauseAudioDevice(cAID, 0)

	debug.LogI("SDL initialized audio device with config: %+v", cfg)
	return nil
}

func (a *audioSystem) update() {
	if a.cAID <= 0 {
		return
	}

	channels := len(a.cfg.Spec.Channels)
	pushSize, track := a.mixer.Mix()
	pushSize *= channels

	for i := 0; i < pushSize; i++ {
		channel := a.cfg.Spec.Channels[i%channels]
		t := track[channel]
		cursor := i / channels
		a.buf[i] = t[cursor]
	}

	if pushSize > 0 {
		C.SDL_QueueAudio(a.cAID, unsafe.Pointer(&a.buf[0]), C.uint(pushSize*int(unsafe.Sizeof(float32(0)))))
	}
}

func (a *audioSystem) destroy() {
	if a.cAID > 0 {
		C.SDL_CloseAudioDevice(a.cAID)
	}

	C.SDL_QuitSubSystem(C.SDL_INIT_AUDIO)
}
