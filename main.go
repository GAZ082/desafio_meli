package main

import (
	meli "github.com/gaz082/desafio_meli/app"
	"github.com/tidwall/gjson"
)

func main() {
	queryData := meli.GetSearchedItemList("samsung", 2, 1)
	itemsData := meli.GetItemData(meli.GetItemIDs(queryData))
	q := gjson.ParseBytes(itemsData)
	meli.WriteCSV("salida.csv", q)
}
