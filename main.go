package main

import (
	"log"

	meli "github.com/gaz082/desafio_meli/app"
	"github.com/tidwall/gjson"
)

func main() {
	query := "samsung"

	queryData := meli.GetSearchedItemList(query, 250, 50)
	ids := meli.GetItemIDs(queryData)

	//api just supports up to 20 items/batch

	log.Printf("ids: %v", len(ids))

	var itemsData [][]byte
	c := 0
	itemBatch := 20
	for i := 0; i <= len(ids); i += itemBatch {
		c++
		log.Printf("tomando desde %v-%v", i, c*itemBatch)
		itemsData = append(itemsData, meli.GetItemData(ids[i:c*itemBatch]))
	}
	var headers bool
	for i, v := range itemsData {
		q := gjson.ParseBytes(v)
		if i == 0 {
			headers = true
		} else {
			headers = false
		}
		meli.WriteCSV("salida.csv", q, query, headers)
	}

}
