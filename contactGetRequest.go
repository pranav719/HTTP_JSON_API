package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	id := "1"
	it := "2020-12-20"
	url := "http://localhost:8080/contacts?user=" + id + "&infection_timestamp=" + it

	//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("If-None-Match", `A/"aaaaa"`)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var d []string

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Fatal(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("The possible victims of Covid are :")
	fmt.Println(d)
}
