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
	else{
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		http.Error(w, "some error in connecting",0)
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		http.Error(w, "some error in ping",0)
		return
	}

	//fmt.Println("Connected to MongoDB!")

	Collection := client.Database("test").Collection("test1")

	switch r.Method {
	case "GET":
		id := strings.SplitN(r.URL.Path, "/", 3)[2]
		d , err:= getUserData(id)
		if err!= nil{
			http.Error(w, "Some Error Occureeeed", 5000)
			break
		}
		fmt.Printf("data for id %d send\n", id)
		/*for k, v := range r.URL.Query() {
			fmt.Printf("%s: %s\n", k, v)
		}*/
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(d)
	case "POST":
		var d User
		/*reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}*/
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, "Some Error Occureeeed", 5000)
			return
		}
		str ,err := insertUserData(d)
		if err != nil{
			http.Error(w, "some error in insertion", 5000)
			return
		}
		fmt.Printf("\n", d)
		//data := "Hello Bhai sab badiya"
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(str)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

if !strings.Contains(r.URL.Path, "/contact") {
		http.NotFound(w, r)
		return
	}
	else{
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		http.Error(w, "some error in connecting",0)
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		http.Error(w, "some error in ping",0)
		return
	}

	//fmt.Println("Connected to MongoDB!")

	Collection2 := client.Database("test").Collection("contactData")

	switch r.Method {
	case "GET":
		/*id := strings.SplitN(r.URL.Path, "/", 3)[2]
		d , err:= getUserData(id)
		if err!= nil{
			http.Error(w, "Some Error Occureeeed", 5000)
			break
		}*/
		var userid string
		var ct time
		for k, v := range r.URL.Query() {
			//fmt.Printf("%s: %s\n", k, v)
			if(k="userid")userid = v
			else if(k ="contacttimestamp") ct=v
			else{
				http.Error(w, "incorrect variable naming", 0)
				return
			}
		}

		contactTime = time.Parse("2020-10-2", ct)
		victims , err:= getContactData(userid, contactTime)
		if err!= nil{
			http.Error(w, "Error in retreving data", 0)
			return
		}


		fmt.Printf("Possible Covid victims can be %+v\n", victims)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(victims)
	case "POST":
		var d contactData
		/*reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}*/
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, "Some Error Occureeeed", 5000)
			log.Fatal(err)
		}
		str ,err := insertContactData(d)
		fmt.Printf("Data received and added in database\n", d)
		//data := "Hello Bhai sab badiya"
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(str)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
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

func getUserData(id string) (User, error){
	filter := bson.D{{"id", id}}

	var result User

	err = Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return User{}, err
	}

	return result, nil	
}

func insertUserData(userData User) (string, error){
	insertResult, err := Collection.InsertOne(context.TODO(), userData)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID, nil
}

func insertContactData(contactData Contact) (string, error){
	insertResult, err := Collection2.InsertOne(context.TODO(), contactData)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID, nil
}

func getContactData(userid string, cTime time) ([]string, error){

	filter := bson.D{{
		"contacttimestamp": bson.D{{
			"$gt": cTime.Add(-14*time.Day),
			"$lt": cTime
		}},
		userid : bson.D{{
			"$in": {"userid1", "userid2"}
		}}
	}}

	var results []string

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := Collection2.Find(context.TODO(), filter)
	if err != nil {
		return [], err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			return [], err
		}

		if(elem.User1ID == userid) results = append(results, elem.User2ID)
		else results = append(results, elem.User1ID)

	}

	return results, nil
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

type Contact struct {
	User1ID                string    `json:"user1id"`
	User2ID                string    `json:"user2id"`
	ContactTimestamp time.Time `json:"contacttimestamp"`
}
