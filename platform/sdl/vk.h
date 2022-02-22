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

#pragma once

#define VK_NO_PROTOTYPES
#include <vulkan/vulkan.h>

#include <SDL2/SDL.h>
#include <SDL2/SDL_vulkan.h>

extern VkResult vkVerifyVersion(uint32_t apiVersion);

extern VkResult vkCreateInstance(const VkInstanceCreateInfo* pCreateInfo,
								 const VkAllocationCallbacks* pAllocator,
								 VkInstance* pInstance);

extern void vkDestroySurface(VkInstance instance,
							 VkSurfaceKHR surface,
							 const VkAllocationCallbacks* pAllocator);

extern void vkDestroyInstance(VkInstance instance,
							  const VkAllocationCallbacks* pAllocator);
