package main

import (
	meli "github.com/gaz082/desafio_meli/app"
)

func main() {
	query := []string{"alimento peces", "alimento gatos", "alimento perros"}
	fileName := "salida.csv"
	meli.WriteHeader(fileName)
	meli.LoadDataToFIle(query, fileName)

}
