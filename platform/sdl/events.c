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

#include <SDL2/SDL.h>
#include "event.h"

void processWindowEvent(goEvent* ge, SDL_WindowEvent e) {
	switch (e.event) {
		case SDL_WINDOWEVENT_NONE:
			ge->windowState = (ge->windowState | WINDOW_CREATED);
			break;
		case SDL_WINDOWEVENT_SHOWN:
		case SDL_WINDOWEVENT_RESTORED:
			ge->windowState = (ge->windowState | WINDOW_SHOWN) & ~WINDOW_HIDDEN;
			break;
		case SDL_WINDOWEVENT_HIDDEN:
		case SDL_WINDOWEVENT_MINIMIZED:
			ge->windowState = (ge->windowState | WINDOW_HIDDEN) & ~WINDOW_SHOWN;
			break;

		case SDL_WINDOWEVENT_MOVED:
		case SDL_WINDOWEVENT_RESIZED:
			ge->windowState = (ge->windowState | WINDOW_RECT_CHANGED);
			break;

		case SDL_WINDOWEVENT_ENTER:
			ge->windowState = (ge->windowState | WINDOW_ENTER) & ~WINDOW_LEAVE;
			break;
		case SDL_WINDOWEVENT_LEAVE:
			ge->windowState = (ge->windowState | WINDOW_LEAVE) & ~WINDOW_ENTER;
			break;

		case SDL_WINDOWEVENT_FOCUS_GAINED:
			ge->windowState =
				(ge->windowState | WINDOW_FOCUS_GAINED) & ~WINDOW_FOCUS_LOST;
			break;
		case SDL_WINDOWEVENT_FOCUS_LOST:
			ge->windowState =
				(ge->windowState | WINDOW_FOCUS_LOST) & ~WINDOW_FOCUS_GAINED;
			break;

		case SDL_WINDOWEVENT_CLOSE:
			ge->windowState = ge->windowState | WINDOW_CLOSE;
			break;
	}
}

int processEvents(goEvent* ge) {
	int alive = 1;
	SDL_Event e;

	while (SDL_PollEvent(&e) != 0) {
		if (e.type == SDL_WINDOWEVENT) {
			if (e.window.windowID == ge->window &&
				!(ge->windowState & WINDOW_CLOSE)) {
				processWindowEvent(ge, e.window);
			}
		}
		if (e.type == SDL_MOUSEWHEEL) {
			ge->mouseWheelX = e.wheel.x;
			ge->mouseWheelY = e.wheel.y;
		}
		if (e.type == SDL_QUIT) {
			// ge->windowState = ge->windowState | WINDOW_CLOSE;
			alive = 0;
		}
	}

	if (!alive) {
		ge->windowState = ge->windowState & ~WINDOW_CLOSE;
	}

	return alive;
}
