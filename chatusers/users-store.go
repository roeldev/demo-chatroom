// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatusers

import (
	"sync"

	"github.com/go-pogo/errors"
	"github.com/google/uuid"
)

const (
	ErrUserNotFound      errors.Msg = "user does not exist"
	ErrNameAlreadyExists errors.Msg = "name is already in use"
)

type UsersStore interface {
	All() map[UserID]User
	Has(id UserID) bool
	Get(id UserID) (User, error)
	Add(user User) (UserID, error)
	Update(id UserID, user User) error
	Delete(id UserID) (User, bool)
}

type usersStore struct {
	mut   sync.RWMutex
	users map[UserID]User
}

func NewUsersStore(size uint64) UsersStore {
	return &usersStore{
		users: make(map[UserID]User, size),
	}
}

func (us *usersStore) All() map[UserID]User {
	us.mut.RLock()
	defer us.mut.RUnlock()

	clone := make(map[UserID]User, len(us.users))
	for k, v := range us.users {
		clone[k] = v
	}
	return clone
}

func (us *usersStore) Has(id UserID) bool {
	_, err := us.Get(id)
	return err == nil
}

func (us *usersStore) Get(id UserID) (User, error) {
	us.mut.RLock()
	defer us.mut.RUnlock()

	u, ok := us.users[id]
	if !ok {
		return u, errors.New(ErrUserNotFound)
	}

	return u, nil
}

func (us *usersStore) Add(user User) (UserID, error) {
	us.mut.Lock()
	defer us.mut.Unlock()

	for _, u := range us.users {
		if u.Name == user.Name {
			return uuid.Nil, errors.New(ErrNameAlreadyExists)
		}
	}

	id := uuid.New()
	us.users[id] = user
	return id, nil
}

func (us *usersStore) Update(id UserID, user User) error {
	us.mut.Lock()
	defer us.mut.Unlock()

	_, ok := us.users[id]
	if !ok {
		return errors.New(ErrUserNotFound)
	}

	us.users[id] = user
	return nil
}

func (us *usersStore) Delete(id UserID) (User, bool) {
	us.mut.Lock()
	defer us.mut.Unlock()

	user, ok := us.users[id]
	if ok {
		delete(us.users, id)
	}
	return user, ok
}
