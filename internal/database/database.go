package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DataBase struct {
	Name string
	Mux  *sync.RWMutex
}

type DBData struct {
	Chirps map[int]Chirp `json:"chirps"`
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
		log.Fatal(err)
	}
}

func (db *DataBase) Read() (DBData, error) {
	db.Mux.Lock()
	defer db.Mux.Unlock()

	data, err := os.ReadFile(db.Name)
	if err != nil {
		log.Printf("error reading file")
		log.Fatal(err)
		return DBData{}, err
	}

	dbData := DBData{}
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		log.Printf("error unmarshaling file")
		log.Fatal(err)
		return DBData{}, err
	}

	return dbData, nil
}

func (db *DataBase) GetChirps() ([]Chirp, error) {
	dbData, err := db.Read()
	if err != nil {
		log.Printf("Cannot get chirps: %s", err)
		return []Chirp{}, err
	}

	chirps := make([]Chirp, 0)
	for _, chirp := range dbData.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DataBase) CreateChirp(body string) (Chirp, error) {
	chirps, err := db.Read()
	if err != nil {
		log.Fatal(err)
		return Chirp{}, err
	}

	newId := len(chirps.Chirps) + 1
	newChirp := Chirp{
		Id:   newId,
		Body: body,
	}

	chirps.Chirps[newId] = newChirp
	db.Write(chirps)
	return newChirp, nil
}

func (db *DataBase) CreateIfNotExits() error {
	fileInfo, err := os.Stat(db.Name)
	if err == nil {
		db.Name = fileInfo.Name()
		return nil
	}

	file, err := os.Create(db.Name)
	if err != nil {
		return err
	}
	db.Name = file.Name()

	err = os.WriteFile(db.Name, []byte("{\"chirps\": {}}"), 0666)
	if err != nil {
		log.Printf("Error initializing file: %s", err)
		return err
	}

	return nil
}
