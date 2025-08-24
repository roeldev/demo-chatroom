// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatusers

import (
	"image/color"
	"math/rand/v2"
)

func RandomColors() (color.RGBA, color.RGBA) {
	c1 := color.RGBA{
		R: uint8(rand.UintN(120)),
		G: uint8(rand.UintN(120)),
		B: uint8(rand.UintN(120)),
	}
	return c1, SecondaryColor(c1)
}

func SecondaryColor(c1 color.Color) color.RGBA {
	//r, g, b, _ := c1.RGBA()
	return color.RGBA{R: 255, G: 255, B: 255}
}
