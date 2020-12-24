# goARRG - go Assembly RequiRed Game-engine
[![Go Reference](https://pkg.go.dev/badge/goarrg.com.svg)](https://pkg.go.dev/goarrg.com)<br/>
goARRG is a assembly required game engine where the pieces may or may not be provided and it is up to you to put them together.
This allows the user to customise the engine almost exactly to their needs.
It is not a goal however to make it easy to rip out a piece and replace it later.

## Development Roadmap
*Not in actual development order
 - ~~MVP for engine development~~
	 - ~~Platform initialization and setup~~
	 - ~~Maths~~
	 - ~~Glue APIs~~
	 - ~~Basic Tooling~~
 - Renderer
 - ECS
 - Editor

## OS Specific Dependencies

| OS | Dependencies |
| -- | -- |
| Ubuntu | sudo apt-get install build-essential libglu1-mesa-dev mesa-common-dev libxext-dev |
| Windows | mingw-w64, cmake |

## Install instructions

Install Go manually to ensure you have the latest version.<br/>
https://golang.org/doc/install

If you want to use vulkan you need the vulkan SDK. Set the env var `VULKAN_SDK` if it is a manual install<br/>
https://vulkan.lunarg.com/sdk/home

If you want to use SDL audio you need the relevant dev library on linux<br/>
This is usually pulseaudio `libpulse-dev` or alsa `libasound2-dev`

<details>
<summary>Manual Install</summary><br>
NOTE: Path must be <code>$HOME/go/...</code> <br/>
Replace <code>$HOME</code> with <code>%USERPROFILE%</code> on windows <br/><br/>

<pre>
go env -w GO111MODULE=off
git clone &ltURL&gt $HOME/go/src/goarrg.com
cd $HOME/go/src/goarrg.com
go get -d ./...
go run goarrg.com/cmd/goarrg build yourself -vv
</pre>
</details>

<details>
<summary>Modules Install</summary><br>
NOTE: If you are switching from manual install, you need to run <code>go env -w GO111MODULE=on</code><br/><br/>

<pre>
mkdir projectfolder
cd projectfolder
go mod init github.com/username/projectname
go get goarrg.com
go run goarrg.com/cmd/goarrg build yourself -vv
</pre>
</details>

## Examples

To run the examples (excluding the shared folder), simply cd to the folder and run `goarrg run -vv`.<br/>
