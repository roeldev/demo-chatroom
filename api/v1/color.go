// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import (
	"encoding/hex"
	"image/color"
	"strings"

	"github.com/go-pogo/errors"
)

const (
	ErrColorTooShort errors.Msg = "value is too short"
)

type ColorDecodeError struct {
	cause    error
	HexValue string
}

func (e *ColorDecodeError) Unwrap() error { return e.cause }

func (e *ColorDecodeError) Error() string {
	return "unable to decode hex color value `" + e.HexValue + "`"
}

// NewColor creates a new [Color] from a [color.Color] by encoding it using
// [Color.Encode].
func NewColor(val color.Color) *Color {
	if val == nil {
		return nil
	}

	x := new(Color)
	x.Encode(val)
	return x
}

// Encode color.Color into Color's Value.
func (x *Color) Encode(c color.Color) {
	r, g, b, _ := color.RGBAModel.Convert(c).RGBA()
	x.Value = "#" + hex.EncodeToString([]byte{uint8(r), uint8(g), uint8(b)})
}

// Decode the Color's Value into a color.Color.
func (x *Color) Decode() (color.Color, error) {
	if x == nil {
		return nil, nil
	}

	v := strings.TrimPrefix(x.Value, "#")
	if len(v) == 3 {
		v = string([]byte{
			v[0], v[0],
			v[1], v[1],
			v[2], v[2],
		})
	}
	if len(v) < 6 {
		return nil, errors.WithStack(&ColorDecodeError{
			cause:    ErrColorTooShort,
			HexValue: x.Value,
		})
	}

	b, err := hex.DecodeString(v)
	if err != nil {
		return nil, errors.WithStack(&ColorDecodeError{
			cause:    err,
			HexValue: x.Value,
		})
	}

	return color.RGBA{R: b[0], G: b[1], B: b[2]}, nil
}
