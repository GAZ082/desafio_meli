/*
1) Barrer una lista de más de 150 ítems ids en el servicio público:

https://api.mercadolibre.com/sites/MLA/search?q=chromecast&limit=50#json

2) Por cada resultado, realizar el correspondiente GET por Item_Id al recurso público:

https://api.mercadolibre.com/items/{Item_Id}
*/

package desafio_meli

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type meliOutput struct {
	id             string
	price          string
	original_price string
	seller_id      string
	permalink      string
	thumbnail      string
}

func GetSearchedItemList(searchQuery string, records, pageSize int) [][]byte {

	//max limit = 50

	pages := int(math.Ceil(float64(records) / float64(pageSize)))
	var outputSlice [][]byte
	for i := 0; i < pages; i++ {
		query := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?q=%v&limit=%v&offset=%v#json",
			searchQuery,
			strconv.Itoa(pageSize),
			strconv.Itoa(i*pageSize))
		outputSlice = append(outputSlice, doQueryReturnData(query))
	}
	return outputSlice
}

func GetItemIDs(input [][]byte) (out []string) {

	for _, d := range input {
		v := gjson.GetBytes(d, "results.#.id")
		for _, item := range v.Array() {
			out = append(out, item.String())
		}
	}

	return
}

func GetItemData(IDSlice []string) []byte {
	strBuild := strings.Builder{}
	for i, s := range IDSlice {
		strBuild.WriteString(s)
		if i < len(IDSlice)-1 {
			strBuild.WriteString(",")
		}
	}
	// query := fmt.Sprintf("https://api.mercadolibre.com/items/%v",

	query := fmt.Sprintf("https://api.mercadolibre.com/items?ids=%v",
		strBuild.String())
	log.Printf("%v", query)
	strBuild.Reset()
	return doQueryReturnData(query)
}

func ParseItemData(input []byte) string {
	return gjson.GetBytes(input, "").String()
}

func WriteCSV(fileName string, input *[]meliOutput) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	for _, record := range *input {
		var row []string
		row = append(row, record.id, record.price, record.original_price,
			record.seller_id, record.permalink, record.thumbnail)
		err = writer.Write(row)
		if err != nil {
			fmt.Println(err)
		}
	}
	writer.Flush()
}

func doQueryReturnData(query string) []byte {
	resp, err := http.Get(query)
	if err != nil {
		log.Println(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Query: %v - %v bytes", query, len(body))
	return body
}
