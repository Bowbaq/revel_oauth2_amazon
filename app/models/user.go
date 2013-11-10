package models

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"sync"
)

type User struct {
	Uid   uint64
	Token oauth.Token
}

var db = struct {
	users   map[uint64]*User
	next_id uint64
	sync.RWMutex
}{
	users: make(map[uint64]*User),
}

func GetUser(id uint64) *User {
	db.RLock()
	user := db.users[id]
	db.RUnlock()

	return user
}

func NewUser() *User {
	db.Lock()
	db.next_id += 1
	user := &User{Uid: db.next_id}
	db.users[user.Uid] = user
	db.Unlock()

	return user
}

func SetToken(id uint64, token *oauth.Token) error {
	user := GetUser(id)
	if user == nil {
		return fmt.Errorf("Couldn't find user with id: %d", id)
	}

	db.Lock()
	db.users[user.Uid].Token = *token
	db.Unlock()

	return nil
}
