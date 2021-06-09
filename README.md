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

## Supported Platforms
Currently goarrg only supports Ubuntu and Windows, 386 and amd64. Vulkan is only supported on amd64.
However, there is nothing preventing you from creating a platform package to support other platforms.

### Cross Compile
There is cross compile support for the supported platforms, assuming you installed a C/C++ cross compiler with the correct file names. To cross compile to other platforms, or to use a non default toolchain, you need to set the `CC`/`CXX`/`AR` environmental variables. For Windows, you also need to set `RC`.

#### Default Compiler Selection
| Taraget Platform | Prefix |
| -- | -- |
| linux_386 | i686-linux-gnu |
| linux_amd64 | x86_64-linux-gnu |
| windows_386 | i686-w64-mingw32 |
| windows_amd64 | x86_64-w64-mingw32 |

`CC={{.Prefix}}-gcc`<br>
`CXX={{.Prefix}}-g++`<br>
`AR={{.Prefix}}-gcc-ar`<br>
**Windows Only:**<br>
`RC={{.Prefix}}-windres`

## Dependencies

goarrg requires go 1.16+
| OS | Dependencies |
| -- | -- |
| Ubuntu | sudo apt-get install build-essential cmake libxext-dev libpulse-dev |
| Windows | mingw-w64, cmake |

### Graphics API Specific
| OS | Folder Prefix | Dependencies |
| -- | -- | -- |
| Ubuntu | gl | sudo apt-get install libglu1-mesa-dev mesa-common-dev |
| Ubuntu_amd64 | vk | Vulkan SDK |
| Windows_amd64 | vk | Vulkan SDK |

## Install instructions

Install Go manually to ensure you have the latest version.<br/>
https://golang.org/doc/install

## Quick Start
<pre><code>mkdir projectfolder
cd projectfolder
go mod init github.com/username/projectname
echo -e "//+build tools\npackage main\nimport _ \"goarrg.com/cmd/goarrg\"" > tools.go
go get -d goarrg.com/...
go run goarrg.com/cmd/goarrg build yourself -vv
</code></pre>

## Examples
Examples can be found at https://github.com/goarrg/examples
