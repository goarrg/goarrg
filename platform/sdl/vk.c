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

#include "vk.h"

#if defined(_WIN32)
#define VKAPI_PTR __stdcall
#else
#define VKAPI_PTR
#endif

typedef struct VkAllocationCallbacks {
} VkAllocationCallbacks;

typedef void(VKAPI_PTR* PFN_vkVoidFunction)(void);
typedef PFN_vkVoidFunction(VKAPI_PTR* PFN_vkGetInstanceProcAddr)(
	VkInstance instance,
	const char* pName);
typedef VkResult(VKAPI_PTR* PFN_vkCreateInstance)(
	const VkInstanceCreateInfo* pCreateInfo,
	const VkAllocationCallbacks* pAllocator,
	VkInstance* pInstance);
typedef void(VKAPI_PTR* PFN_vkDestroySurfaceKHR)(
	VkInstance instance,
	VkSurfaceKHR surface,
	const VkAllocationCallbacks* pAllocator);
typedef void(VKAPI_PTR* PFN_vkDestroyInstance)(
	VkInstance instance,
	const VkAllocationCallbacks* pAllocator);

VkResult vkCreateInstance(VkApplicationInfo appInfo,
						  VkInstanceCreateInfo createInfo,
						  VkInstance* pInstance) {
	PFN_vkGetInstanceProcAddr vkGetInstanceProcAddr =
		(PFN_vkGetInstanceProcAddr)SDL_Vulkan_GetVkGetInstanceProcAddr();

	if (!vkGetInstanceProcAddr) {
		return VK_ERROR_INVALID_EXTERNAL_HANDLE;
	}

	PFN_vkCreateInstance vkCreateInstance =
		(PFN_vkCreateInstance)vkGetInstanceProcAddr(NULL, "vkCreateInstance");

	if (!vkCreateInstance) {
		return VK_ERROR_INVALID_EXTERNAL_HANDLE;
	}

	createInfo.pApplicationInfo = &appInfo;
	return vkCreateInstance(&createInfo, NULL, pInstance);
}

void vkDestroySurface(VkInstance instance, VkSurfaceKHR surface) {
	PFN_vkGetInstanceProcAddr vkGetInstanceProcAddr =
		(PFN_vkGetInstanceProcAddr)SDL_Vulkan_GetVkGetInstanceProcAddr();

	if (!vkGetInstanceProcAddr) {
		return;
	}

	PFN_vkDestroySurfaceKHR vkDestroySurfaceKHR =
		(PFN_vkDestroySurfaceKHR)vkGetInstanceProcAddr(NULL,
													   "vkDestroySurfaceKHR");

	if (!vkDestroySurfaceKHR) {
		return;
	}

	vkDestroySurfaceKHR(instance, surface, NULL);
}

void vkDestroyInstance(VkInstance instance) {
	PFN_vkGetInstanceProcAddr vkGetInstanceProcAddr =
		(PFN_vkGetInstanceProcAddr)SDL_Vulkan_GetVkGetInstanceProcAddr();

	if (!vkGetInstanceProcAddr) {
		return;
	}

	PFN_vkDestroyInstance vkDestroyInstance =
		(PFN_vkDestroyInstance)vkGetInstanceProcAddr(NULL, "vkDestroyInstance");

	if (!vkDestroyInstance) {
		return;
	}

	vkDestroyInstance(instance, NULL);
}
