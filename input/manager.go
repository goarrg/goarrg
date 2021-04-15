/*
Copyright 2021 The goARRG Authors.

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

package input

import (
	"sync"
)

var manager struct {
	sync.Map
	sync.Mutex
}

/*
	RegisterDevice registers a device that can be obetained with the GetDevice*
	functions. The instance must be valid even if the device disconnected and
	shall receive future events when the device reconnected.

	Whether the instance refers to the same physical device on reconnect is a
	implementation detail.

	Unplugging Controller A and plugging in Controller B of the same type, may
	send events to the instance of the unplugged Controller A even tho they are
	not the same physical device. However if the player then plugs Controller A,
	it must use another unplugged instance or register a new instance.
*/
func RegisterDevice(d Device) {
	manager.Lock()
	defer manager.Unlock()

	v, ok := manager.LoadOrStore(d.Type(), []Device{d})

	if ok {
		l := v.([]Device)
		manager.Store(d.Type(), append(l, d))
	}
}

/*
	DeviceOfType retruns the first found device of the given type or nil.
*/
func DeviceOfType(t string) Device {
	v, ok := manager.Load(t)

	if ok {
		return v.([]Device)[0]
	}

	return nil
}

/*
	DevicesOfType retruns a copy of the slice containing all the devices of
	a given type or nil.
*/
func DevicesOfType(t string) []Device {
	v, ok := manager.Load(t)

	if ok {
		return append([]Device(nil), v.([]Device)...)
	}

	return nil
}

/*
	ScanMask represents a bitmask of the types of input events to scan for.
*/
type ScanMask uint8

const (
	ScanValue = 1 << iota
	ScanAxis
	ScanCoords
)

/*
	Scan returns the device and the action, using mask to filter action types,
	that had a DeviceAction triggered this frame or nil and 0.

	This is useful for input mapping without having to specifically code to
	support every device type.

	As this is only meant to be used for key mapping, we can assume there will
	only be one input a frame and that speed isn't too important so just check
	everything.
*/
func Scan(mask ScanMask) (device Device, action DeviceAction) {
	if mask == 0 {
		return nil, 0
	}

	manager.Range(func(key, value interface{}) bool {
		for _, d := range value.([]Device) {
			for i := DeviceAction(0); i < (^DeviceAction(0)); i++ {
				if d.ActionStartedFor(i) {
					switch d.StateFor(i).(type) {
					case Value:
						if mask&ScanValue == 0 {
							continue
						}
					case Axis:
						if mask&ScanAxis == 0 {
							continue
						}
					case Coords:
						if mask&ScanCoords == 0 {
							continue
						}
					}
					device = d
					action = i
					return false
				}
			}
		}
		return true
	})

	return device, action
}
