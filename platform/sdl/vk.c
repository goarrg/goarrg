//go:build !disable_vk && amd64
// +build !disable_vk,amd64

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

VkResult vkCreateInstance(const VkInstanceCreateInfo* pCreateInfo,
						  const VkAllocationCallbacks* pAllocator,
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

	return vkCreateInstance(pCreateInfo, pAllocator, pInstance);
}

void vkDestroySurface(VkInstance instance,
					  VkSurfaceKHR surface,
					  const VkAllocationCallbacks* pAllocator) {
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

	vkDestroySurfaceKHR(instance, surface, pAllocator);
}

void vkDestroyInstance(VkInstance instance,
					   const VkAllocationCallbacks* pAllocator) {
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

	vkDestroyInstance(instance, pAllocator);
}
