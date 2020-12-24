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

type Renderer interface {
	Draw() float64
	Resize(int, int)
	Shutdown()
	Destroy()
}

type GLProfile uint8

const (
	GLProfileCore GLProfile = iota
	GLProfileCompat
	GLProfileES
)

type GLConfig struct {
	Major   uint8
	Minor   uint8
	Profile GLProfile
}

type GLInstance interface {
	ProcAddr() uintptr
	SwapBuffers()
}

type GLRenderer interface {
	Renderer
	GLConfig() GLConfig
	GLInit(GLInstance) error
}

type VkConfig struct {
	API        uint32
	Layers     []string
	Extensions []string
}

type VkInstance interface {
	Ptr() uintptr
	ProcAddr() uintptr
	Surface() uintptr
}

type VkRenderer interface {
	Renderer
	VkConfig() VkConfig
	VkInit(VkInstance) error
}
