package reader

import (
	"encoding/json"
	stockModel "homework/loms/internal/model/stock"
	"io"
	"log"
	"os"
)

func ReadStocks(filePath string) []stockModel.Stock {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("File open error: %v", err.Error())
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("File read error: %v", err.Error())
	}

	var stocks []stockModel.Stock
	err = json.Unmarshal(byteValue, &stocks)
	if err != nil {
		log.Fatalf("JSON parsing error: %v", err.Error())
	}

	return stocks
}
