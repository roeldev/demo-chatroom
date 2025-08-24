// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor_Decode(t *testing.T) {
	tests := map[string]struct {
		wantColor color.Color
		wantErr   error
	}{
		"#000": {
			wantColor: color.RGBA{},
		},
		"fF0000": {
			wantColor: color.RGBA{R: 255},
		},
	}

	for str, tc := range tests {
		t.Run(str, func(t *testing.T) {
			haveColor, haveErr := (&Color{Value: str}).Decode()
			if tc.wantErr == nil {
				assert.NoError(t, haveErr)
			} else {
				assert.ErrorIs(t, haveErr, tc.wantErr)
			}

			assert.Equal(t, tc.wantColor, haveColor)
		})
	}
}
