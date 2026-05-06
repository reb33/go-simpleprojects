package main

import (
	"3-struct/bins"
	"3-struct/storage"
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {
	loadAndPrintBins()
}

func createAndSaveBins() {
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

func loadAndPrintBins() {
	binList, err := storage.LoadBins()
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, bin := range binList.Bins {
		fmt.Println(bin)
	}
}
