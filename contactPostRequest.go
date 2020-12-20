package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Contact struct {
	User1ID          string    `json:"user1id"`
	User2ID          string    `json:"user2id"`
	ContactTimestamp time.Time `json:"contacttimestamp"`
}

func main() {
	url := "http://localhost:8080/contacts"

	ct := time.Now().Add(-13 * 24 * time.Hour)
	fmt.Println(ct)
	contact := Contact{
		User1ID:          "1",
		User2ID:          "4",
		ContactTimestamp: ct,
	}
	fmt.Println("data send : \n", contact)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(contact)

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
