package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"work/internal/browser"
)

func main() {
	ctx, allocCxl, ctxCxl := browser.GetContext()
	defer allocCxl()
	defer ctxCxl()

	connStr := "host=postgres port=5432 user=sample password=sample dbname=sample sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	age := 21
	rows, err := db.Query("SELECT $1", age)
	column := 0
	for rows.Next() {
		err = rows.Scan(&column)
	}
	if err != nil {
		log.Fatal(err)
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
