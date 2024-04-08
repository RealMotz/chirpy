package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type DataBase struct {
	Name string
	Mux  *sync.RWMutex
}

type DBData struct {
	Chirps map[int]Chirp        `json:"chirps"`
	Users  map[int]User         `json:"users"`
	Tokens map[string]time.Time `json:"tokens"`
}

func (db *DataBase) Write(data DBData) {
	db.Mux.Lock()
	defer db.Mux.Unlock()

	encodedData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error encoding data: %s", err)
		return
	}
	err = os.WriteFile(db.Name, encodedData, 0666)
	if err != nil {
		log.Printf("error writing file")
		log.Fatal(err)
	}
}

func (db *DataBase) Read() (DBData, error) {
	db.CreateIfNotExits()

	db.Mux.Lock()
	defer db.Mux.Unlock()

	data, err := os.ReadFile(db.Name)
	if err != nil {
		log.Printf("error reading file")
		return DBData{}, err
	}

	dbData := DBData{}
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		log.Printf("error unmarshaling file")
		return DBData{}, err
	}

	return dbData, nil
}

func (db *DataBase) CreateIfNotExits() {
	_, err := os.Stat(db.Name)
	if err == nil {
		return
	}

	dbStructure := DBData{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
		Tokens: map[string]time.Time{},
	}
	db.Write(dbStructure)
}
