package main

import (
	"log"

	"work/internal/browser"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type TPage struct {
	gorm.Model
	Name        string `gorm:"size:255"`
	URL         string `gorm:"size:255"`
	Description string `gorm:"size:255"`
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
	db.Select("name, url, description").Find(&TPages)
	for _, v := range TPages {
		log.Printf("selected %s (%s): '%s'", v.Name, v.URL, v.Description)
	}

	db.Exec("TRUNCATE TABLE t_pages")
}
