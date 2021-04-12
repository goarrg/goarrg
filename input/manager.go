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

import "sync"

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
	GetDeviceOfType retruns the first found device of the given type or nil.
*/
func GetDeviceOfType(t string) Device {
	v, ok := manager.Load(t)

	if ok {
		return v.([]Device)[0]
	}

	return nil
}

/*
	GetDevicesOfType retruns a copy of the slice containing all the devices of
	a given type or nil.
*/
func GetDevicesOfType(t string) []Device {
	v, ok := manager.Load(t)

	if ok {
		return append([]Device(nil), v.([]Device)...)
	}

	return nil
}
