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

package goarrg

import (
	"goarrg.com/asset/audio"
)

/*
Refer to driver's documentation on the valid value ranges
*/
type AudioConfig struct {
	Spec audio.Spec

	// Length of the buffer in terms of samples, this is usually a power of 2
	BufferLength int
}

type Audio interface {
	AudioConfig() AudioConfig
	Init(AudioConfig) error
	Mix() (int, audio.Track)
	Update()
	Destroy()
}

type audioNull struct{}

func (audioNull) AudioConfig() AudioConfig { return AudioConfig{} }
func (audioNull) Init(AudioConfig) error   { return nil }
func (audioNull) Mix() (int, audio.Track)  { return 0, nil }
func (audioNull) Update()                  {}
func (audioNull) Destroy()                 {}
