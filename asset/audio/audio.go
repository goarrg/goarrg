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

package audio

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"goarrg.com/asset"
	"goarrg.com/debug"
)

type Channel uint32

const (
	ChannelLeft Channel = iota
	ChannelRight
	ChannelCenter
	ChannelLowFrequency
	ChannelSurroundLeft
	ChannelSurroundRight
	ChannelBackSurroundLeft
	ChannelBackSurroundRight
	ChannelCount = iota
)

/*
	Track in 32 bit Float non-interleaved format, map key represents the individual channels

	Driver is then responsible for converting it into the appropriate signal for output
	audio is non-interleaved to make it easier to support different channel orders
*/
type Track map[Channel][]float32

type Spec struct {
	Channels  []Channel
	Frequency int
}

type Asset interface {
	Track() Track
	Spec() Spec
	DurationSeconds() float64
	DurationSamples() int
}

type assetImpl struct {
	spec            Spec
	durationSeconds float64
	durationSamples int
	track           Track
}

type format struct {
	magic  []byte
	decode func(asset.Asset) (Spec, int, []float32, error)
}

var mtx = sync.Mutex{}
var formats = atomic.Value{}

func ChannelsMono() []Channel {
	return []Channel{ChannelLeft}
}

func ChannelsStereo() []Channel {
	return []Channel{ChannelLeft, ChannelRight}
}

func Channels5Point1() []Channel {
	return []Channel{ChannelLeft, ChannelRight, ChannelCenter, ChannelLowFrequency, ChannelSurroundLeft, ChannelSurroundRight}
}

func Channels7Point1() []Channel {
	return []Channel{ChannelLeft, ChannelRight, ChannelCenter, ChannelLowFrequency, ChannelBackSurroundLeft, ChannelBackSurroundRight, ChannelSurroundLeft, ChannelSurroundRight}
}

func RegisterFormat(magic string, decode func(asset.Asset) (Spec, int, []float32, error)) {
	mtx.Lock()
	f, _ := formats.Load().([]format)
	formats.Store(append(f, format{[]byte(magic), decode}))
	mtx.Unlock()
}

func Load(file string) (Asset, error) {
	a, err := asset.Load(file)

	if err != nil {
		return nil, debug.ErrorWrap(err, "Failed to load audio")
	}

	r := bufio.NewReader(a.Reader())
	formats, _ := formats.Load().([]format)

formats:
	for _, f := range formats {
		if a.Size() < len(f.magic) {
			continue
		}

		magic, err := r.Peek(len(f.magic))

		if err != nil && err != io.EOF {
			return nil, debug.ErrorWrap(err, "Failed to load audio")
		}

		for i, b := range f.magic {
			if b != magic[i] && b != '?' {
				continue formats
			}
		}

		spec, samples, interleavedTrack, err := f.decode(a)

		if err != nil {
			return nil, err
		}

		duration := float64(samples) / float64(spec.Frequency) / float64(len(spec.Channels))
		track := make(Track)

		for i, s := range interleavedTrack {
			track[spec.Channels[i%len(spec.Channels)]] = append(track[spec.Channels[i%len(spec.Channels)]], s)
		}

		return &assetImpl{
			spec,
			duration,
			samples / len(spec.Channels),
			track,
		}, nil
	}

	return nil, debug.ErrorNew("Failed to load audio, unknown format")
}

func (s *assetImpl) Track() Track {
	return s.track
}

func (s *assetImpl) Spec() Spec {
	return s.spec
}

func (s *assetImpl) DurationSeconds() float64 {
	return s.durationSeconds
}

func (s *assetImpl) DurationSamples() int {
	return s.durationSamples
}

func (c Channel) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fallthrough
	case 's':
		switch c {
		case ChannelLeft:
			fmt.Fprintf(s, "L")
		case ChannelRight:
			fmt.Fprintf(s, "R")
		case ChannelCenter:
			fmt.Fprintf(s, "C")
		case ChannelLowFrequency:
			fmt.Fprintf(s, "LF")
		case ChannelSurroundLeft:
			fmt.Fprintf(s, "SL")
		case ChannelSurroundRight:
			fmt.Fprintf(s, "SR")
		case ChannelBackSurroundLeft:
			fmt.Fprintf(s, "BSL")
		case ChannelBackSurroundRight:
			fmt.Fprintf(s, "BSR")
		default:
			fmt.Fprintf(s, "User%d", uint32(c)-ChannelCount+1)
		}
	}
}
