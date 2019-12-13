// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in13

// epd2in13 module should be connected in following pins
const (
	PinRST = "RST" // 17 for Rpi
	PinDC  = "DC"  // 25 for Rpi
	// PinCS   = "CS"   // 8 for Rpi spi0.0
	PinBusy = "BUSY" // 24 for Rpi
)

const (
	epd2in13Width  = 122
	epd2in13Height = 250
)
