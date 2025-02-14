//go:build !goarrg_disable_vk
// +build !goarrg_disable_vk

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

	#include "vk.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	"goarrg.com"
	"goarrg.com/debug"
)

type vkInstance struct {
	procAddr  uintptr
	cInstance C.VkInstance
	cSurface  C.VkSurfaceKHR
}

func (vk *vkInstance) destroy() {
	if vk.cSurface != C.VK_NULL_HANDLE {
		C.vkDestroySurface(vk.cInstance, vk.cSurface)
		vk.cSurface = C.VK_NULL_HANDLE
	}
	if vk.cInstance != nil {
		C.vkDestroyInstance(vk.cInstance)
		vk.cInstance = nil
	}
}

func (vk *vkInstance) ProcAddr() uintptr {
	return vk.procAddr
}

func (vk *vkInstance) Uintptr() uintptr {
	return uintptr(unsafe.Pointer(vk.cInstance))
}

func (vk *vkInstance) CreateSurface() (uint64, error) {
	if vk.cSurface != C.VK_NULL_HANDLE {
		C.vkDestroySurface(vk.cInstance, vk.cSurface)
	}
	//nolint:staticcheck
	if C.SDL_Vulkan_CreateSurface(Platform.display.mainWindow.cWindow, vk.cInstance, &vk.cSurface) != C.SDL_TRUE {
		err := debug.Errorf("Failed to create surface: %s", C.GoString(C.SDL_GetError()))
		C.SDL_ClearError()
		return 0, err
	}
	return uint64(vk.cSurface), nil
}

func vkResultStr(code C.VkResult) string {
	switch code {
	case C.VK_SUCCESS:
		return "VK_SUCCESS"
	case C.VK_NOT_READY:
		return "VK_NOT_READY"
	case C.VK_TIMEOUT:
		return "VK_TIMEOUT"
	case C.VK_EVENT_SET:
		return "VK_EVENT_SET"
	case C.VK_EVENT_RESET:
		return "VK_EVENT_RESET"
	case C.VK_INCOMPLETE:
		return "VK_INCOMPLETE"
	case C.VK_ERROR_OUT_OF_HOST_MEMORY:
		return "VK_ERROR_OUT_OF_HOST_MEMORY"
	case C.VK_ERROR_OUT_OF_DEVICE_MEMORY:
		return "VK_ERROR_OUT_OF_DEVICE_MEMORY"
	case C.VK_ERROR_INITIALIZATION_FAILED:
		return "VK_ERROR_INITIALIZATION_FAILED"
	case C.VK_ERROR_DEVICE_LOST:
		return "VK_ERROR_DEVICE_LOST"
	case C.VK_ERROR_MEMORY_MAP_FAILED:
		return "VK_ERROR_MEMORY_MAP_FAILED"
	case C.VK_ERROR_LAYER_NOT_PRESENT:
		return "VK_ERROR_LAYER_NOT_PRESENT"
	case C.VK_ERROR_EXTENSION_NOT_PRESENT:
		return "VK_ERROR_EXTENSION_NOT_PRESENT"
	case C.VK_ERROR_FEATURE_NOT_PRESENT:
		return "VK_ERROR_FEATURE_NOT_PRESENT"
	case C.VK_ERROR_INCOMPATIBLE_DRIVER:
		return "VK_ERROR_INCOMPATIBLE_DRIVER"
	case C.VK_ERROR_TOO_MANY_OBJECTS:
		return "VK_ERROR_TOO_MANY_OBJECTS"
	case C.VK_ERROR_FORMAT_NOT_SUPPORTED:
		return "VK_ERROR_FORMAT_NOT_SUPPORTED"
	case C.VK_ERROR_FRAGMENTED_POOL:
		return "VK_ERROR_FRAGMENTED_POOL"
	case C.VK_ERROR_OUT_OF_POOL_MEMORY:
		return "VK_ERROR_OUT_OF_POOL_MEMORY"
	case C.VK_ERROR_INVALID_EXTERNAL_HANDLE:
		return "VK_ERROR_INVALID_EXTERNAL_HANDLE"
	case C.VK_ERROR_SURFACE_LOST_KHR:
		return "VK_ERROR_SURFACE_LOST_KHR"
	case C.VK_ERROR_NATIVE_WINDOW_IN_USE_KHR:
		return "VK_ERROR_NATIVE_WINDOW_IN_USE_KHR"
	case C.VK_SUBOPTIMAL_KHR:
		return "VK_SUBOPTIMAL_KHR"
	case C.VK_ERROR_OUT_OF_DATE_KHR:
		return "VK_ERROR_OUT_OF_DATE_KHR"
	case C.VK_ERROR_INCOMPATIBLE_DISPLAY_KHR:
		return "VK_ERROR_INCOMPATIBLE_DISPLAY_KHR"
	case C.VK_ERROR_VALIDATION_FAILED_EXT:
		return "VK_ERROR_VALIDATION_FAILED_EXT"
	case C.VK_ERROR_INVALID_SHADER_NV:
		return "VK_ERROR_INVALID_SHADER_NV"
	case C.VK_ERROR_INVALID_DRM_FORMAT_MODIFIER_PLANE_LAYOUT_EXT:
		return "VK_ERROR_INVALID_DRM_FORMAT_MODIFIER_PLANE_LAYOUT_EXT"
	case C.VK_ERROR_FRAGMENTATION_EXT:
		return "VK_ERROR_FRAGMENTATION_EXT"
	case C.VK_ERROR_NOT_PERMITTED_EXT:
		return "VK_ERROR_NOT_PERMITTED_EXT"
	case C.VK_ERROR_INVALID_DEVICE_ADDRESS_EXT:
		return "VK_ERROR_INVALID_DEVICE_ADDRESS_EXT"
	case C.VK_ERROR_FULL_SCREEN_EXCLUSIVE_MODE_LOST_EXT:
		return "VK_ERROR_FULL_SCREEN_EXCLUSIVE_MODE_LOST_EXT"
	}

	if code < 0 {
		return fmt.Sprintf("Unknown VkResult error (%d)", code)
	}

	return fmt.Sprintf("Unknown VkResult (%d)", code)
}

type vkWindow struct {
	cfg        goarrg.VkConfig
	renderer   goarrg.VkRenderer
	vkInstance *vkInstance

	windowW int
	windowH int
}

func vkInit(r goarrg.VkRenderer) error {
	Platform.logger.IPrintf("Creating vk Window")

	if r == nil {
		err := debug.Errorf("Invalid renderer")
		Platform.logger.EPrintf("failed to create window: Invalid renderer")
		return err
	}

	err := createWindow(C.SDL_WINDOW_VULKAN)
	if err != nil {
		return err
	}

	defer func() {
		if Platform.display.mainWindow != nil && Platform.display.mainWindow.api == nil {
			Platform.display.mainWindow.destroy()
			Platform.display.mainWindow = nil
		}
	}()

	vkCfg := r.VkConfig()
	cInstance := C.VkInstance(nil)

	Platform.logger.IPrintf("Renderer requested config: %+v", vkCfg)

	{
		cNumSDLExt := C.uint(0)
		if C.SDL_Vulkan_GetInstanceExtensions(Platform.display.mainWindow.cWindow, &cNumSDLExt, nil) != C.SDL_TRUE {
			err := debug.Errorf("Failed to get list of SDL required vulkan extensions: %s", C.GoString(C.SDL_GetError()))
			C.SDL_ClearError()
			return err
		}

		cExt := (**C.char)(C.calloc((C.size_t(cNumSDLExt) + C.size_t(len(vkCfg.Extensions))), C.size_t(unsafe.Sizeof((*C.char)(nil)))))
		defer C.free(unsafe.Pointer(cExt))
		ext := unsafe.Slice((**C.char)(unsafe.Pointer(cExt)), int(cNumSDLExt)+len(vkCfg.Extensions))

		if C.SDL_Vulkan_GetInstanceExtensions(Platform.display.mainWindow.cWindow, &cNumSDLExt, cExt) != C.SDL_TRUE {
			err := debug.Errorf("Failed to get list of SDL required vulkan extensions: %s", C.GoString(C.SDL_GetError()))
			C.SDL_ClearError()
			return err
		}

		for i, e := range vkCfg.Extensions {
			ext[int(cNumSDLExt)+i] = C.CString(e)
			defer C.free(unsafe.Pointer(ext[int(cNumSDLExt)+i]))
		}

		cLayers := (**C.char)(C.calloc(C.size_t(len(vkCfg.Layers)), C.size_t(unsafe.Sizeof((*C.char)(nil)))))
		defer C.free(unsafe.Pointer(cLayers))
		layers := unsafe.Slice((**C.char)(unsafe.Pointer(cLayers)), len(vkCfg.Layers))

		for i, l := range vkCfg.Layers {
			layers[i] = C.CString(l)
			defer C.free(unsafe.Pointer(layers[i]))
		}

		cVkAInfo := C.VkApplicationInfo{
			sType: C.VK_STRUCTURE_TYPE_APPLICATION_INFO,
		}

		if vkCfg.API != 0 {
			cVkAInfo.apiVersion = C.uint32_t(vkCfg.API)
		} else {
			cVkAInfo.apiVersion = C.VK_API_VERSION_1_0
		}

		cVkCInfo := C.VkInstanceCreateInfo{
			sType:                   C.VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO,
			enabledLayerCount:       C.uint32_t(len(vkCfg.Layers)),
			ppEnabledLayerNames:     cLayers,
			enabledExtensionCount:   cNumSDLExt + C.uint32_t(len(vkCfg.Extensions)),
			ppEnabledExtensionNames: cExt,
		}

		//nolint:staticcheck
		if ret := C.vkCreateInstance(cVkAInfo, cVkCInfo, &cInstance); ret != C.VK_SUCCESS {
			for i := 0; i < int(cNumSDLExt); i++ {
				vkCfg.Extensions = append(vkCfg.Extensions, C.GoString(ext[i]))
			}

			if ret == C.VK_ERROR_INVALID_EXTERNAL_HANDLE {
				return debug.ErrorWrapf(debug.Errorf("Failed to find vulkan loader"), "Failed to create vk instance with config %#v", vkCfg)
			}

			return debug.ErrorWrapf(debug.Errorf("%s", vkResultStr(ret)), "Failed to create vk instance with config %#v", vkCfg)
		}
	}

	window := &vkWindow{
		cfg:      vkCfg,
		renderer: r,
		vkInstance: &vkInstance{
			procAddr:  uintptr(C.SDL_Vulkan_GetVkGetInstanceProcAddr()),
			cInstance: cInstance,
		},
	}

	if err := r.VkInit(window.vkInstance); err != nil {
		window.vkInstance.destroy()
		return err
	}

	Platform.display.mainWindow.api = window
	Platform.display.mainWindow.api.resize(Platform.config.Window.Rect.W, Platform.config.Window.Rect.H)
	Platform.logger.IPrintf("Created vk window")

	return nil
}

func (vkw *vkWindow) resize(w int, h int) {
	vkw.windowW = w
	vkw.windowH = h

	if w != 0 && h != 0 {
		var cW, cH C.int
		C.SDL_Vulkan_GetDrawableSize(Platform.display.mainWindow.cWindow, &cW, &cH)
		vkw.renderer.Resize(int(cW), int(cH))
	} else {
		vkw.renderer.Resize(0, 0)
	}
}

func (vkw *vkWindow) destroy() {
	vkw.vkInstance.destroy()
}
