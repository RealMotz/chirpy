package database

import "log"

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DataBase) CreateUser(email string) (User, error) {
	dbData, err := db.Read()
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}

	newId := len(dbData.Users) + 1
	user := User{
		Id:    newId,
		Email: email,
	}

	dbData.Users[newId] = user
	db.Write(dbData)
	return user, nil
}
