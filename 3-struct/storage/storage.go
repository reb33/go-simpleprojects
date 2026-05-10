package storage

import (
	"3-struct/bins"
	"encoding/json"
)

type DB interface {
	Read() ([]byte, error)
	Write(data []byte) error
}

type Storage struct {
	DB
}

func (storage *Storage) LoadBins() (*bins.BinList, error){
	data, error := storage.Read()
	if error != nil {
		return nil, error
	}
	var bins bins.BinList
	error = json.Unmarshal(data, &bins)
	if error != nil {
		return nil, error
	}
	return &bins, nil
}

func (storage *Storage) SaveBins(bins *bins.BinList) error {
	json, error := json.Marshal(bins)
	if error != nil {
		return error
	}
	return storage.Write(json)
}