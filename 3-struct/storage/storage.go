package storage

import (
	"3-struct/bins"
	"encoding/json"
	"errors"
	"os"
)

type DB interface {
	Read() ([]byte, error)
	Write(data []byte) error
}

type Storage struct {
	DB
}

func NewStorage(db DB) *Storage {
	storage := Storage{DB: db}
	_, err := storage.DB.Read()
	if err != nil {
		if errors.Is(err, os.ErrNotExist){
			storage.SaveBins(bins.NewBinList(nil))
		}
	}
	return &storage
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

func (storage *Storage) AddBin(bin *bins.Bin) error {
	bins, err := storage.LoadBins()
	if err != nil {
		return err
	}
	bins.AddBin(bin)
	return storage.SaveBins(bins)
}

func (storage *Storage) DelBinById(binId string) error {
	bins, err := storage.LoadBins()
	if err != nil {
		return err
	}
	bins.DelBinById(binId)
	return storage.SaveBins(bins)
}