// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatusers

import (
	"image/color"
	"strings"

	"github.com/go-pogo/errors"
	"github.com/google/uuid"
)

type UserID = uuid.UUID

func IdentifierString(uid UserID, user UserDetails) string {
	return uid.String() + " (" + user.Name + ")"
}

type Flag uint8

func (f Flag) Has(flag Flag) bool { return f&flag != 0 }

//goland:noinspection GoSnakeCaseUsage
const (
	Flag_None Flag = iota
	Flag_IsBot
)

type Status uint8

//goland:noinspection GoSnakeCaseUsage
const (
	Status_Default Status = iota
	Status_Busy
	Status_Away
)

type User struct {
	UserDetails
	Flags  Flag
	Status Status
}

type UserDetails struct {
	Name     string
	Initials string
	Color1   color.Color // primary color
	Color2   color.Color // secondary color
}

const (
	ErrEmptyName            errors.Msg = "user name should not be empty"
	ErrInvalidNameBot       errors.Msg = "only bots may use a 'bot' suffix"
	ErrUnsufficientContrast errors.Msg = "unsufficient color contrast"
)

type UserOption func(u *User) error

// NewUser creates a new valid User
func NewUser(name string, opts ...UserOption) (*User, error) {
	user := User{
		UserDetails: UserDetails{
			Name: name,
		},
	}

	var err error
	for _, opt := range opts {
		err = errors.Append(err, opt(&user))
	}

	if user.Name == "" {
		return nil, errors.New(ErrEmptyName)
	} else if strings.HasSuffix(user.Name, "bot") && !user.Flags.Has(Flag_IsBot) {
		return nil, errors.New(ErrInvalidNameBot)
	}

	if user.Initials == "" {
		user.Initials = InitialsFromName(user.Name)
	}

	// make sure colors are set
	if user.Color1 == nil {
		user.Color1, user.Color2 = RandomColors()
	} else if user.Color2 == nil {
		user.Color2 = SecondaryColor(user.Color1)
		//} else {
		// check color contrast?
	}

	return &user, nil
}

func InitialsFromName(name string) string {
	if name == "" {
		return ""
	}

	parts := strings.Fields(name)
	if len(parts) >= 2 {
		return strings.ToUpper(string(parts[0][0]) + string(parts[1][0]))
	}
	return strings.ToUpper(string(parts[0][0]))
}

func (u User) IsZero() bool { return u == User{} }

func (ud UserDetails) IsZero() bool { return ud == UserDetails{} }
