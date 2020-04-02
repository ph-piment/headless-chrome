package main

import (
	"fmt"
	"log"

	"work/internal/browser"

	"work/pkg/redis"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// TPage is row for t_pages table.
type TPage struct {
	gorm.Model
	Name        string `gorm:"column:name" "size:255"`
	URL         string `gorm:"column:url" "size:255"`
	Description string `gorm:"column:description" "size:255"`
}

func main() {
	ctx, allocCxl, ctxCxl := browser.GetContext()
	defer allocCxl()
	defer ctxCxl()

	db, err := gorm.Open("postgres", "host=postgres port=5432 user=sample dbname=sample password=sample sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	URL := "https://github.com/avelino/awesome-go"
	sect := "Selenium and browser control tools."
	res, err := browser.ListAwesomeGoProjects(ctx, URL, sect)
	if err != nil {
		log.Fatalf("could not list awesome go projects: %v", err)
	}

	var page TPage
	for k, v := range res {
		log.Printf("register %s (%s): '%s'", k, v.URL, v.Description)
		page = TPage{Name: k, URL: v.URL, Description: v.Description}
		db.Create(&page)
	}

	var TPages []TPage
	db.Find(&TPages)
	for k, v := range TPages {
		log.Printf("selected %d %d %s %s %s %s %s %s", k, v.ID, v.Name, v.URL, v.Description, v.CreatedAt, v.UpdatedAt, v.DeletedAt)
	}

	db.First(&page)
	page.Name = "aaa"
	page.URL = "bbb"
	page.Description = "ccc"
	db.Save(&page)

	db.Find(&TPages)
	for k, v := range TPages {
		log.Printf("selected2 %d %d %s %s %s %s %s %s", k, v.ID, v.Name, v.URL, v.Description, v.CreatedAt, v.UpdatedAt, v.DeletedAt)
	}

	db.Exec("TRUNCATE TABLE t_pages")

	fmt.Println("---------------------------------------------------")
	fmt.Println("Start StringGetSet")
	stringKey := "StringGetSet_Key"
	stringValue := "StringGetSet_Val"
	redis.SetString(stringKey, stringValue)
	redis.SetStringWithExpire(stringKey, stringValue, 0)
	redis.GetString(stringKey)
	fmt.Println("---------------------------------------------------")
	fmt.Println("Start ListGetSet")
	listKey := "ListGetSet_Key"
	listValue := []string{"val1", "va2", "val3"}
	redis.RPush(listKey, listValue)
	redis.LPush(listKey, listValue)
	redis.LSet(listKey, 0, "valval")
	redis.LRange(listKey, 1, 5)
	redis.AllRange(listKey)
	redis.LIndex(listKey, 6)
	fmt.Println("---------------------------------------------------")
	fmt.Println("Start HashGetSet")
	hashKey := "HashGetSet_Key"
	for field, val := range map[string]string{"field1": "val1", "field2": "val2"} {
		err := redis.HSet(hashKey, field, val)
		if err != nil {
			fmt.Println("redis.Client.HSet Error:", err)
		}
	}
	redis.HGet(hashKey, "field2")
	redis.HGetAll(hashKey)
	fmt.Println("---------------------------------------------------")
}
