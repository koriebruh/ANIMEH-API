package main

import (
	"fmt"
	"koriebruh/find/conf"
	"koriebruh/find/insert"
	"log"
	"testing"
	"time"
)

func TestInsertDB(t *testing.T) {
	fmt.Println("INSERT MAU JALAN NIH")
	time.Sleep(25 * time.Second)
	log.Println("INSERT UNIT TEST REK")
	config := conf.GetConfig()
	fmt.Println(config.Elastic)
	ES := conf.ElasticClient(config)
	insert.InsertDataCSVToElastic(ES)
}
