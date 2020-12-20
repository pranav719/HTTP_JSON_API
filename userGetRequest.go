package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Dob               string    `json:"dob"`
	PhoneNo           int       `json:"phoneno"`
	Email             string    `json:"email"`
	CreationTimestamp time.Time `json:"creationtimestamp"`
}

func main() {

	id := "4"
	url := "http://localhost:8080/users/" + id

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

	var d User

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Fatal(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("The data for the user with id %d is :", id)
	fmt.Println(d)
}
