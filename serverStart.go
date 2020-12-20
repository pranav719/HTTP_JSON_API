package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func receive(w http.ResponseWriter, r *http.Request) {
	//timestamp := time.Now()

	if !strings.Contains(r.URL.Path, "/users") {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case "GET":
		id := strings.SplitN(r.URL.Path, "/", 3)[2]
		fmt.Printf("%T\n", id)
		/*for k, v := range r.URL.Query() {
			fmt.Printf("%s: %s\n", k, v)
		}*/
		w.Write([]byte("Received a GET request\n"))
	case "POST":
		var d User
		/*reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}*/
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, "Some Error Occureeeed", 5000)
			log.Fatal(err)
		}

		fmt.Printf("\n", d)
		//data := "Hello Bhai sab badiya"
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(d)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}

	/*reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	//w.Write([]byte("Received a POST request\n"))

	/*err := json.Unmarshal(jsonStr, &data)
	    if err != nil {
	        fmt.Println(err)
		}*/

	// validation
	// if error in input return unsuccessful

	// put data into database

	//retrun successful

}

func test() {
	fmt.Println("hello")
}

func main() {
	http.HandleFunc("/users", receive)
	http.HandleFunc("/users/", receive)
	test()
	http.ListenAndServe(":8080", nil)

}

type User struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Dob               string    `json:"dob"`
	PhoneNo           int       `json:"phoneno"`
	Email             string    `json:"email"`
	CreationTimestamp time.Time `json:"creationtimestamp"`
}
