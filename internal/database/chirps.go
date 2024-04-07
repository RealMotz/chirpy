package database

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DataBase) GetChirps() ([]Chirp, error) {
	dbData, err := db.Read()
	if err != nil {
		log.Printf("Cannot read database: %s", err)
		return []Chirp{}, err
	}

	chirps := make([]Chirp, 0)
	for _, chirp := range dbData.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DataBase) GetChirp(id string) (Chirp, error) {
	dbData, err := db.Read()
	if err != nil {
		log.Printf("Cannot read database: %s", err)
		return Chirp{}, err
	}
	chirpId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Id not valid: %s", err)
		return Chirp{}, err
	}
	if chirp, ok := dbData.Chirps[chirpId]; ok {
		return chirp, nil
	}

	return Chirp{}, errors.New("Chirp not found")
}

func (db *DataBase) CreateChirp(body string) (Chirp, error) {
	dbData, err := db.Read()
	if err != nil {
		fmt.Printf("error creating chirp: %s", err)
		return Chirp{}, err
	}

	newId := len(dbData.Chirps) + 1
	chirp := Chirp{
		Id:   newId,
		Body: body,
	}

	dbData.Chirps[newId] = chirp
	db.Write(dbData)
	return chirp, nil
}
