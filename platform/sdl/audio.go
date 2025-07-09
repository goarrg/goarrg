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
	#cgo pkg-config: sdl3
	#include <SDL3/SDL.h>

	static SDL_IOStream* IOFromUintptr(uintptr_t ptr, int sz) {
		return SDL_IOFromConstMem((void*)ptr, sz);
	}
*/
import "C"

import (
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
	cfg     goarrg.AudioConfig
	mixer   goarrg.Audio
	cStream *C.SDL_AudioStream
	buf     []float32
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
		return nil, debug.Errorf("Unsupported channel value %d", n)
	}
}

func verifyChannelList(list []audio.Channel) ([]audio.Channel, error) {
	m := audioChannelMask(0)

	for _, c := range list {
		m |= audioChannelMask(1 << c)
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
		return nil, debug.Errorf("Unsupported channel list %v", list)
	}
}

func decodeWAV(a *asset.File) (audio.Spec, int, []float32, error) {
	cIO := C.IOFromUintptr(C.uintptr_t(a.Uintptr()), C.int(a.Size()))
	cSpec := C.SDL_AudioSpec{}
	cBuf := (*C.Uint8)(nil)
	cLen := C.Uint32(0)

	//nolint:staticcheck
	if !C.SDL_LoadWAV_IO(cIO, true, &cSpec, &cBuf, &cLen) {
		err := debug.ErrorWrapf(debug.Errorf("%s", C.GoString(C.SDL_GetError())), "Failed to decode WAV")
		C.SDL_ClearError()
		return audio.Spec{}, 0, nil, err
	}
	defer C.SDL_free(unsafe.Pointer(cBuf))

	channels, err := channelCountToList(int(cSpec.channels))
	if err != nil {
		return audio.Spec{}, 0, nil, debug.ErrorWrapf(err, "Failed to decode WAV")
	}

	cStream := C.SDL_CreateAudioStream(&cSpec, &C.SDL_AudioSpec{
		freq:     cSpec.freq,
		format:   C.SDL_AUDIO_F32,
		channels: cSpec.channels,
	})
	if cStream == nil {
		err := debug.ErrorWrapf(debug.Errorf("%s", C.GoString(C.SDL_GetError())), "Failed to decode WAV")
		C.SDL_ClearError()
		return audio.Spec{}, 0, nil, err
	}
	defer C.SDL_DestroyAudioStream(cStream)

	if !C.SDL_PutAudioStreamData(cStream, unsafe.Pointer(cBuf), C.int(cLen)) {
		err := debug.ErrorWrapf(debug.Errorf("%s", C.GoString(C.SDL_GetError())), "Failed to decode WAV")
		C.SDL_ClearError()
		return audio.Spec{}, 0, nil, err
	}
	if !C.SDL_FlushAudioStream(cStream) {
		err := debug.ErrorWrapf(debug.Errorf("%s", C.GoString(C.SDL_GetError())), "Failed to decode WAV")
		C.SDL_ClearError()
		return audio.Spec{}, 0, nil, err
	}
	samples := int(cLen) / (int(cSpec.format&0xFF) / 8)
	track := make([]float32, samples)
	if C.SDL_GetAudioStreamData(cStream, unsafe.Pointer(unsafe.SliceData(track)), C.int(unsafe.Sizeof(float32(0)))*C.int(samples)) == -1 {
		err := debug.ErrorWrapf(debug.Errorf("%s", C.GoString(C.SDL_GetError())), "Failed to decode WAV")
		C.SDL_ClearError()
		return audio.Spec{}, 0, nil, err
	}

	return audio.Spec{
		Frequency: int(cSpec.freq),
		Channels:  channels,
	}, samples, track, nil
}

func (*platform) AudioInit(mixer goarrg.Audio) error {
	if Platform.config.Audio.Importer.EnableWAV {
		audio.RegisterFormat("RIFF????WAVE", decodeWAV)
	}

	if mixer == nil {
		Platform.logger.IPrintf("SDL audio disabled")
		return nil
	}

	if !C.SDL_InitSubSystem(C.SDL_INIT_AUDIO) {
		err := debug.ErrorWrapf(debug.Errorf("%s", C.GoString(C.SDL_GetError())), "Failed to init SDL audio")
		C.SDL_ClearError()
		return err
	}

	cfg := mixer.AudioConfig()
	if len(cfg.Spec.Channels) == 0 {
		return debug.ErrorWrapf(debug.Errorf("No channels defined"), "Failed to init SDL audio")
	}

	channels, err := verifyChannelList(cfg.Spec.Channels)
	if err != nil {
		return debug.ErrorWrapf(err, "Failed to init SDL audio")
	}

	spec := C.SDL_AudioSpec{
		freq:     C.int(cfg.Spec.Frequency),
		format:   C.SDL_AUDIO_F32,
		channels: C.int(len(channels)),
	}

	//nolint:staticcheck
	cStream := C.SDL_OpenAudioDeviceStream(C.SDL_AUDIO_DEVICE_DEFAULT_PLAYBACK, &spec, nil, nil)
	if cStream == nil {
		err := debug.ErrorWrapf(debug.Errorf("%s", C.GoString(C.SDL_GetError())), "Failed to init SDL audio")
		C.SDL_ClearError()
		return err
	}
	if err := mixer.Init(platformInterface{}, cfg); err != nil {
		return debug.ErrorWrapf(err, "Failed to init SDL audio")
	}
	C.SDL_ResumeAudioDevice(C.SDL_GetAudioStreamDevice(cStream))

	Platform.audio.cStream = cStream
	Platform.audio.mixer = mixer
	Platform.audio.cfg = cfg
	Platform.audio.buf = make([]float32, cfg.Spec.Frequency*len(cfg.Spec.Channels))

	Platform.logger.IPrintf("Initialized audio device with config: %+v", cfg)
	return nil
}

func (a *audioSystem) update() {
	if a.cStream == nil {
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
		C.SDL_PutAudioStreamData(a.cStream, unsafe.Pointer(unsafe.SliceData(a.buf)), C.int(pushSize*int(unsafe.Sizeof(float32(0)))))
	}
}

func (a *audioSystem) destroy() {
	if a.cStream != nil {
		C.SDL_DestroyAudioStream(a.cStream)
	}

	C.SDL_QuitSubSystem(C.SDL_INIT_AUDIO)
}
