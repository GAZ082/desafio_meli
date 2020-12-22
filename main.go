package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"

	meli "github.com/gaz082/desafio_meli/app"
	"github.com/tidwall/gjson"
)

func main() {
	query := []string{"alimento peces", "alimento gatos", "alimento perros"}
	fileName := "salida.csv"
	header(fileName)
	for _, search := range query {
		queryData := meli.GetSearchedItemList(url.QueryEscape(search), 1000, 50)
		ids := meli.GetItemIDs(queryData)

		log.Printf("ids: %v", len(ids))

		var itemsData [][]byte
		const itemBatch = 20 //max 20
		cut := len(ids) / itemBatch
		var from, to int
		for i := 0; i <= cut; i++ {
			from = i * itemBatch
			if i != cut {
				to = i*itemBatch + itemBatch
			} else {
				to = i*itemBatch + len(ids) - i*itemBatch
			}
			log.Printf("tomando desde %v-%v", from, to)
			itemsData = append(itemsData, meli.GetItemData(ids[from:to]))
		}

		for _, v := range itemsData {
			q := gjson.ParseBytes(v)
			meli.WriteCSV(fileName, q, search)
		}
	}

}

func header(fileName string) {

	csvFile, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	err = writer.Write(meli.Headers)
	if err != nil {
		fmt.Println(err)
	}
	writer.Flush()
}
