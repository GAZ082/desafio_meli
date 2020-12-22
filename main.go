package main

import (
	"log"

	meli "github.com/gaz082/meli/app"
)

func main() {
	queryData := meli.GetSearchedItemList("samsung", 20, 10)

	itemsID := meli.GetItemIDs(queryData)

	log.Printf("%v", itemsID)

}
