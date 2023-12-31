package storage

import (
	"os"
	"time"
)

var Store = `users.json`

type User struct {
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email,omitempty"`
}

type UserList map[string]User

type UserStore struct {
	Increment int      `json:"increment"`
	List      UserList `json:"list"`
}

type UsersFile struct {
	FileJSON     *os.File
	UsersStorage UserStore
}
