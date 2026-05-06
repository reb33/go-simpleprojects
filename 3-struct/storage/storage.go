package storage

import (
	"3-struct/bins"
	"3-struct/files"
	"encoding/json"
)

const FILE_NAME = "bins.json"

func LoadBins() (*bins.BinList, error){
	data, error := files.ReadFile(FILE_NAME)
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

func SaveBins(bins *bins.BinList) error {
	json, error := json.Marshal(bins)
	if error != nil {
		return error
	}
	return files.WriteFile(FILE_NAME, json)
}