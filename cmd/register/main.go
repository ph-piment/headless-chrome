package main

import (
	"log"

	"work/internal/browser"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	ctx, allocCxl, ctxCxl := browser.GetContext()
	defer allocCxl()
	defer ctxCxl()

	db, err := gorm.Open("postgres", "host=postgres port=5432 user=sample dbname=sample password=sample sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	age := 21
	rows, err := db.Raw("SELECT ? AS num", age).Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	column := 0
	for rows.Next() {
		rows.Scan(&column)
	}

	log.Printf("select :%v", column)

	URL := "https://github.com/avelino/awesome-go"
	sect := "Selenium and browser control tools."
	res, err := browser.ListAwesomeGoProjects(ctx, URL, sect)
	if err != nil {
		log.Fatalf("could not list awesome go projects: %v", err)
	}

	for k, v := range res {
		log.Printf("project %s (%s): '%s'", k, v.URL, v.Description)
	}
}
