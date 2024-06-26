package database

import (
	"errors"
	"log"
	"sort"
	"strconv"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

var ErrorChirpNotFound = errors.New("chirp not found")

func (db *DataBase) GetChirps(authorId int, descSort bool) ([]Chirp, error) {
	dbData, err := db.Read()
	if err != nil {
		log.Printf("Cannot read database: %s", err)
		return []Chirp{}, err
	}

	chirps := make([]Chirp, 0)
	for _, chirp := range dbData.Chirps {
		if authorId != 0 {
			if chirp.AuthorId == authorId {
				chirps = append(chirps, chirp)
			}
			continue
		}
		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool {
		if descSort {
			return chirps[i].Id > chirps[j].Id
		} else {
			return chirps[i].Id < chirps[j].Id

		}
	})
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

	chirp, ok := dbData.Chirps[chirpId]
	if !ok {
		return Chirp{}, ErrorChirpNotFound
	}

	return chirp, nil
}

func (db *DataBase) CreateChirp(body string, authorId int) (Chirp, error) {
	dbData, err := db.Read()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbData.Chirps) + 1
	chirp := Chirp{
		Id:       newId,
		Body:     body,
		AuthorId: authorId,
	}

	dbData.Chirps[newId] = chirp
	db.Write(dbData)
	return chirp, nil
}

func (db *DataBase) DeleteChirp(id string) error {
	dbData, err := db.Read()
	if err != nil {
		return err
	}

	chirpId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	if _, ok := dbData.Chirps[chirpId]; !ok {
		return ErrorChirpNotFound
	}

	delete(dbData.Chirps, chirpId)
	db.Write(dbData)
	return nil
}
