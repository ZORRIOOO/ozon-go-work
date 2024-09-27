package reader

import (
	"encoding/json"
	"errors"
	"fmt"
	stockModel "homework/loms/internal/model/stock"
	"io"
	"os"
)

func ReadStocks(filePath string) ([]stockModel.Stock, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []stockModel.Stock{}, errors.New(fmt.Sprintf("Error opening file: %s", err.Error()))
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return []stockModel.Stock{}, errors.New(fmt.Sprintf("File read error: %v", err.Error()))
	}

	var stocks []stockModel.Stock
	err = json.Unmarshal(byteValue, &stocks)
	if err != nil {
		return []stockModel.Stock{}, errors.New(fmt.Sprintf("JSON parsing error: %v", err.Error()))
	}

	return stocks, nil
}
