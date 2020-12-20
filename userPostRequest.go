package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	url := "http://localhost:8080/users"

	user := User{
		ID:                "4",
		Name:              "Ashish",
		Dob:               "2000-10-02",
		PhoneNo:           7849322792,
		Email:             "user4@gmail.com",
		CreationTimestamp: time.Now(),
	}
	fmt.Println("data send : \n", user)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(user)

	//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, payloadBuf)
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
