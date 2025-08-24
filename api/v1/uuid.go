// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import (
	"github.com/google/uuid"
)

// NewUUID creates a new [UUID] from an [uuid.UUID].
func NewUUID(v uuid.UUID) *UUID {
	if v == uuid.Nil {
		return nil
	}

	return &UUID{Value: v.String()}
}

func (x *UUID) ParseUUID() (uuid.UUID, error) {
	if x == nil || x.Value == "" {
		return uuid.Nil, nil
	}
	return uuid.Parse(x.Value)
}
