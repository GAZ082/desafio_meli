package desafio_meli

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
		"body.seller_address.city.name",
		"body.seller_address.state.id",
		"body.seller_address.state.name",
		"body.seller_address.country.id",
		"body.catalog_product_id",
		"body.inventory_id",
		"category_name",
		"search_query",
	}
	Headers = []string{
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
		"seller_address_city_name",
		"seller_address_state_id",
		"seller_address_state_name",
		"seller_address_country_id",
		"catalog_product_id",
		"inventory_id",
		"category_name",
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
	query := fmt.Sprintf("https://api.mercadolibre.com/items?ids=%v",
		strBuild.String())
	strBuild.Reset()
	return doQueryReturnData(query)
}

func ParseItemData(input []byte) string {
	return gjson.GetBytes(input, "").String()
}

func WriteCSV(fileName string, input gjson.Result, search_query string, categories map[string]string) {
	csvFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
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
				continue
			}
			if field == "category_name" {
				cat := record.Get("body.category_id").String()
				row = append(row, categories[cat])
				continue
			}
			row = append(row, s)

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

type categoryRecord struct {
	id           string
	jsonResponse []byte
}

func GetCategories(categoryID []string) map[string]string {
	var temp []categoryRecord
	output := make(map[string]string)
	for _, v := range categoryID {
		q := fmt.Sprintf("https://api.mercadolibre.com/categories/%v", v)
		temp = append(temp, categoryRecord{
			id:           v,
			jsonResponse: doQueryReturnData(q),
		})
	}

	for _, t := range temp {
		d := gjson.GetBytes(t.jsonResponse, "path_from_root.#.name")
		for i, vv := range d.Array() {
			if i == 1 {
				output[t.id] = vv.String()
			}

		}
	}

	return output
}

func WriteHeader(fileName string) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	err = writer.Write(Headers)
	if err != nil {
		fmt.Println(err)
	}
	writer.Flush()
}

func LoadDataToFIle(query []string, fileName string) {
	cats := GetCategories([]string{"MLA1077", "MLA1087", "MLA1094"})
	for _, search := range query {
		queryData := GetSearchedItemList(url.QueryEscape(search), 1000, 50) //1000 50
		ids := GetItemIDs(queryData)
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
			itemsData = append(itemsData, GetItemData(ids[from:to]))
		}
		for _, v := range itemsData {
			q := gjson.ParseBytes(v)
			WriteCSV(fileName, q, search, cats)
		}
	}
}
