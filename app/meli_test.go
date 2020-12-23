package desafio_meli

import (
	"log"
	"testing"
)

func TestGetCategories(t *testing.T) {
	cats := GetCategories([]string{"MLA1077", "MLA1087", "MLA1094"})

	log.Printf("%v", cats["MLA1077"])

}
