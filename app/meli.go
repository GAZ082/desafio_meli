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
	}
)

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
	query := fmt.Sprintf("https://api.mercadolibre.com/items?ids=%v",
		strBuild.String())
	// log.Printf("%v", query)
	strBuild.Reset()
	return doQueryReturnData(query)
}

func ParseItemData(input []byte) string {
	return gjson.GetBytes(input, "").String()
}

func WriteCSV(fileName string, input gjson.Result) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	err = writer.Write(headers)
	if err != nil {
		fmt.Println(err)
	}
	for _, record := range input.Array() {
		var row []string
		for i, field := range fields {
			s := record.Get(field).String()
			if i == 0 && s != "200" { // to filter out invalid records
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
