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

#include <stdint.h>

// clang-format off
#define WINDOW_CREATED      (0x1u << 0)

#define WINDOW_SHOWN        (0x1u << 1)
#define WINDOW_HIDDEN       (0x1u << 2)

#define WINDOW_MOVED        (0x1u << 3)
#define WINDOW_RESIZED      (0x1u << 4)

#define WINDOW_ENTER        (0x1u << 5)
#define WINDOW_LEAVE        (0x1u << 6)

#define WINDOW_FOCUS_GAINED (0x1u << 7)
#define WINDOW_FOCUS_LOST   (0x1u << 8)

#define WINDOW_CLOSE        (0x1u << 9)
// clang-format on

typedef struct {
	uint32_t window;
	uint32_t windowState;
	int32_t windowX;
	int32_t windowY;
	int32_t windowW;
	int32_t windowH;

	int32_t mouseWheelX;
	int32_t mouseWheelY;
} goEvent;
