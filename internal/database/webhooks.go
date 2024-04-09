package database

import "errors"

const (
	UpgradedEvent = "user.upgraded"
)

func (db *DataBase) AddMembership(userId int) error {
	dbData, err := db.Read()
	if err != nil {
		return err
	}

	user, ok := dbData.Users[userId]
	if !ok {
		return errors.New("user not found")
	}

	user.IsChirpRed = true
	dbData.Users[userId] = user
	db.Write(dbData)

	return nil
}
