package main

import (
	"3-struct/api"
	"3-struct/bins"
	"3-struct/config"
	"3-struct/files"
	"3-struct/storage"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

const FILE_NAME = "bins.json"
const API_URL = "https://api.jsonbin.io/v3/b"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("не удалось загрузить .env")
	}
	
	// --- Определяем флаги ---
	create := flag.Bool("create", false, "создать новый bin")
	update := flag.Bool("update", false, "обновить существующий bin")
	get := flag.Bool("get", false, "получить bin по ID")
	delete := flag.Bool("delete", false, "удалить bin по ID")
	list := flag.Bool("list", false, "список всех bins")

	// Дополнительные аргументы
	binID := flag.String("id", "", "ID бина (для get, update, delete)")
	fileName := flag.String("file", "", "путь к файлу (для create, update)")
	binName := flag.String("name", "", "имя бина (по умолчанию — имя файла без расширения)")

	flag.Parse()

	switch {
		case *create:
			if *fileName == "" {
				log.Fatal("--create требует --file")
			}
			createAndSaveBins(getApi(), getStorage(), *fileName, *binName)
		case *get:
			if *binID == "" {
				log.Fatal("--get требует --id")
			}
			getBin(getApi(), *binID)
		case *delete:
			if *binID == "" {
				log.Fatal("--delete требует --id")
			}
			deleteBin(getApi(), getStorage(), *binID)
		case *update:
			if *fileName == "" {
				log.Fatal("--update требует --file")
			}
			if *binID == "" {
				log.Fatal("--update требует --id")
			}
			updateBin(getApi(), *fileName, *binID)
		case *list:
			listBins(getStorage())
		default:
			log.Fatal("требуется одно из действий: --create, --update, --get, --delete, --list")
	}
}


func getApi() *api.Api {
	config := config.NewConfig(API_URL)
	api := api.NewApi(*config)
	return api
}

func getStorage() *storage.Storage {
	fileDB, err := files.NewFileDB(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.NewStorage(fileDB)
	return storage
}

func getFileNameWithoutExt(filePath string) string {
	filename := filepath.Base(filePath)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func createAndSaveBins(api *api.Api, storage *storage.Storage, fileName string, binName string) {
	data, err := files.Read(fileName)
	if err != nil {
		log.Fatal(err)
	}
	metadata, err := api.CreateBin(data)
	if err != nil {
		log.Fatal(err)
	}
	if binName == "" {
		binName = getFileNameWithoutExt(fileName)
	}
	bin, err := bins.NewBin(metadata.Id, metadata.Private, binName, metadata.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
	err = storage.AddBin(bin)
	if err != nil {
		log.Fatal(err)
	}
}

func getBin(api *api.Api, id string) {
	if id == "" {}
	resp, err := api.GetBin(id)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint(resp)
}

func prettyPrint(jsonData []byte) {
    var raw any
    if err := json.Unmarshal(jsonData, &raw); err != nil {
        fmt.Println(string(jsonData))
    }

    pretty, err := json.MarshalIndent(raw, "", "  ")
    if err != nil {
        fmt.Println(string(jsonData))
    }

    fmt.Println(string(pretty))
}

func deleteBin(api *api.Api, storage *storage.Storage, id string) {
	err := api.DeleteBin(id)
	if err != nil {
		log.Fatal(err)
	}
	storage.DelBinById(id)
}

func updateBin(api *api.Api, fileName string, id string) {
	data, err := files.Read(fileName)
	if err != nil {
		log.Fatal(err)
	}
	api.UpdateBin(id, data)
}

func listBins(storage *storage.Storage) {
	binList, err := storage.LoadBins()
	if err != nil {
		log.Fatal(err)
	}
	for _, bin := range binList.Bins {
		fmt.Printf("%s - %s\n", bin.Id, bin.Name)
	}
}

