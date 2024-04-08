package database

import (
	"errors"
	"time"
)

func (db *DataBase) Revoke(token string) error {
	dbData, err := db.Read()
	if err != nil {
		return err
	}

	if _, ok := dbData.Tokens[token]; ok {
		return errors.New("token has already been revoked")
	}

	dbData.Tokens[token] = time.Now()
	db.Write(dbData)
	return nil
}

func (db *DataBase) GetRevokedToken(token string) (time.Time, error) {
	dbData, err := db.Read()
	if err != nil {
		return time.Time{}, err
	}

	revokedTime, ok := dbData.Tokens[token]
	if !ok {
		return time.Time{}, errors.New("no token found")
	}

	return revokedTime, nil
}
