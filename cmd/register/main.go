package main

import (
	"fmt"
	"log"

	"work/internal/browser"

	"github.com/go-redis/redis"
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

	db.Exec("TRUNCATE TABLE t_pages")

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("Redis client:", client)

	fmt.Println("---------------------------------------------------")
	fmt.Println("Start StringGetSet")
	StringGetSet(client)
	fmt.Println("---------------------------------------------------")
	fmt.Println("Start ListGetSet")
	ListGetSet(client)
	fmt.Println("---------------------------------------------------")
	fmt.Println("Start HashGetSet")
	HashGetSet(client)
	fmt.Println("---------------------------------------------------")
}

func StringGetSet(client *redis.Client) {
	key := "StringGetSet_Key"
	// Set
	err := client.Set(key, "StringGetSet_Val", 0).Err()
	if err != nil {
		fmt.Println("redis.Client.Set Error:", err)
	}

	// Get
	val, err := client.Get(key).Result()
	if err != nil {
		fmt.Println("redis.Client.Get Error:", err)
	}
	fmt.Println(val)
}

func ListGetSet(client *redis.Client) {
	key := "ListGetSet_Key"
	// Set
	listVal := []string{"val1", "va2", "val3"}
	err := client.RPush(key, listVal).Err()
	if err != nil {
		fmt.Println("redis.Client.RPush Error:", err)
	}

	// Get
	// Get by lrange
	lrangeVal, err := client.LRange(key, 0, -1).Result()
	if err != nil {
		fmt.Println("redis.Client.LRange Error:", err)
	}
	fmt.Println(lrangeVal)
	// Get by lindex
	lindexVal, err := client.LIndex(key, 2).Result()
	if err != nil {
		fmt.Println("redis.Client.LIndex Error:", err)
	}
	fmt.Println(lindexVal)
}

func HashGetSet(client *redis.Client) {
	key := "HashGetSet_Key"
	// Set
	for field, val := range map[string]string{"field1": "val1", "field2": "val2"} {
		fmt.Println("Inserting", "field:", field, "val:", val)
		err := client.HSet(key, field, val).Err()
		if err != nil {
			fmt.Println("redis.Client.HSet Error:", err)
		}
	}

	// Get
	// HGet(key, field string) *StringCmd
	hgetVal, err := client.HGet(key, "field1").Result()
	if err != nil {
		fmt.Println("redis.Client.HGet Error:", err)
	}
	fmt.Println(hgetVal)

	// HGetAll
	hgetallVal, err := client.HGetAll(key).Result()
	if err != nil {
		fmt.Println("redis.Client.HGetAll Error:", err)
	}
	// fmt.Println("reflect.TypeOf(hgetallVal):", reflect.TypeOf(hgetallVal)) // map[string]string
	fmt.Println(hgetallVal)
}
