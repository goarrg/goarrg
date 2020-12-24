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

package gmath

func insertSpace1u32(x uint16) uint32 {
	m := uint32(x)
	m = (m ^ (m << 8)) & 0x00ff00ff
	m = (m ^ (m << 4)) & 0x0f0f0f0f
	m = (m ^ (m << 2)) & 0x33333333
	m = (m ^ (m << 1)) & 0x55555555
	return m
}

func insertSpace2u32(x uint16) uint32 {
	m := uint32(x)
	m &= 0x000003ff
	m = (m ^ (m << 16)) & 0xff0000ff
	m = (m ^ (m << 8)) & 0x0300f00f
	m = (m ^ (m << 4)) & 0x030c30c3
	m = (m ^ (m << 2)) & 0x09249249
	return m
}

func Morton2u32(x, y uint16) uint32 {
	return insertSpace1u32(x) | insertSpace1u32(y)<<1
}

// 10 bits per value so max value is 1023
func Morton3u32(x, y, z uint16) uint32 {
	return insertSpace2u32(x) | insertSpace2u32(y)<<1 | insertSpace2u32(z)<<2
}
