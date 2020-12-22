package main

import (
	"log"

	meli "github.com/gaz082/desafio_meli/app"
	"github.com/tidwall/gjson"
)

func main() {
	queryData := meli.GetSearchedItemList("samsung", 2, 1)
	itemsData := meli.GetItemData(meli.GetItemIDs(queryData))

	q := gjson.ParseBytes(itemsData)

	for _, v := range q.Array() {
		log.Printf("%v", v.Get("code").String())

	}

	// gjson.ForEachLine(gjson.ParseBytes(itemsData).String(),
	// 	func(line gjson.Result) bool {
	// 		println(line.Get("#id").String())
	// 		return true
	// 	})

}
