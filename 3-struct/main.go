package main

import (
	"3-struct/bins"
	"3-struct/files"
	"3-struct/storage"
	"fmt"
	"log"
	"strconv"
	"time"
)

const FILE_NAME = "bins.json"

func main() {
	fileDB, err := files.NewFileDB(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.Storage{
		DB: fileDB,
	}
	createAndSaveBins(&storage)
	loadAndPrintBins(&storage)
}

func createAndSaveBins(storage *storage.Storage) {
	_bins := []*bins.Bin{}
	for i := range(3){
		bin, err := bins.NewBin(strconv.Itoa(i+1), true, "Bin "+strconv.Itoa(i+1), time.Now())
		if err != nil {
			log.Fatal(err)
			return
		}
		_bins = append(_bins, bin)
	}
	binList := bins.NewBinList(_bins)
	
	err := storage.SaveBins(binList)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func loadAndPrintBins(storage *storage.Storage) {
	binList, err := storage.LoadBins()
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, bin := range binList.Bins {
		fmt.Println(bin)
	}
}
