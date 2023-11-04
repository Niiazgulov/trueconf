package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

/*
Функция для создания нового объекта структуры UsersFile.
Объект хранит информацию об открытом файле со всеми данными о пользователях.
*/
func NewFileStorage() (UsersFile, error) {
	f, err := os.OpenFile(Store, os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		return UsersFile{}, fmt.Errorf("GetRepository: unable to open file: %w", err)
	}
	s := UserStore{}
	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return UsersFile{}, fmt.Errorf("NewFileStorage: unable to decode file: %w", err)
	}
	return UsersFile{
		FileJSON:     f,
		UsersStorage: s,
	}, nil
}
