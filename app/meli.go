package desafio_meli

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

var (
	fields = []string{
		"code",
		"body.id",
		"body.title",
		"body.seller_id",
		"body.category_id",
		"body.price",
		"body.base_price",
		"body.original_price",
		"body.currency_id",
		"body.initial_quantity",
		"body.available_quantity",
		"body.sold_quantity",
		"body.pictures.#(url).secure_url",
		"body.seller_address.city.name",
		"body.seller_address.state.id",
		"body.seller_address.country.id",
		"body.catalog_product_id",
		"body.domain_id",
		"search_query",
	}
	headers = []string{
		"code",
		"id",
		"title",
		"seller_id",
		"category_id",
		"price",
		"base_price",
		"original_price",
		"currency_id",
		"initial_quantity",
		"available_quantity",
		"sold_quantity",
		"picture_url",
		"seller_address_city_name",
		"seller_address_state_id",
		"seller_address_country_id",
		"catalog_product_id",
		"domain_id",
		"search_query",
	}
)

func GetSearchedItemList(searchQuery string, records, pageSize int) [][]byte {

	//max limit = 50

	pages := int(float64(records) / float64(pageSize))
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
	for cd, d := range input {
		log.Printf("BATCH %v ==============", cd)
		v := gjson.GetBytes(d, "results.#.id")
		log.Printf("v.Array() %v", len(v.Array()))
		for ci, item := range v.Array() {
			log.Printf("%v", item.String())
			if ci%10 == 0 && ci > 0 {
				log.Println("")
			}
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
	query := fmt.Sprintf("https://api.mercadolibre.com/items?ids=%v",
		strBuild.String())
	strBuild.Reset()
	return doQueryReturnData(query)
}

func ParseItemData(input []byte) string {
	return gjson.GetBytes(input, "").String()
}

func WriteCSV(fileName string, input gjson.Result, search_query string, header bool) {
	csvFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	if header {
		err = writer.Write(headers)
	}
	if err != nil {
		fmt.Println(err)
	}
	for _, record := range input.Array() {
		var row []string
		for i, field := range fields {
			s := record.Get(field).String()
			if i == 0 && s != "200" { // to filter out invalid records
				log.Println("No OK!")
				continue
			}
			if field == "search_query" {
				row = append(row, search_query)
			} else {
				row = append(row, s)
			}
		}
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
