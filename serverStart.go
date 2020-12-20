package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func receive(w http.ResponseWriter, r *http.Request) {
	//timestamp := time.Now()

	if strings.Contains(r.URL.Path, "/users") {

		switch r.Method {
		case "GET":
			id := strings.SplitN(r.URL.Path, "/", 3)[2]
			d, err := getUserData(id)
			if err != nil {
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
			str, err := insertUserData(d)
			if err != nil {
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
	} else if strings.Contains(r.URL.Path, "/contact") {

		switch r.Method {
		case "GET":
			/*id := strings.SplitN(r.URL.Path, "/", 3)[2]
			d , err:= getUserData(id)
			if err!= nil{
			http.Error(w, "Some Error Occureeeed", 5000)
			break
			}*/
			var userid string
			var contactTime time.Time
			for k, v := range r.URL.Query() {
				//fmt.Printf("%s: %s\n", k, v)
				if k == "userid" {
					userid = v[0]
				} else if k == "contacttimestamp" {
					contactTime, _ = time.Parse("2020-10-2", v[0])
				} else {
					http.Error(w, "incorrect variable naming", 0)
					return
				}
			}

			//contactTime, _ := time.Parse("2020-10-2", ct)
			victims, err := getContactData(userid, contactTime)
			if err != nil {
				http.Error(w, "Error in retreving data", 0)
				return
			}

			fmt.Printf("Possible Covid victims can be %+v\n", victims)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(victims)

		case "POST":
			var d Contact
			/*reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}*/
			if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
				http.Error(w, "Some Error Occureeeed", 5000)
				return
			}
			str, err := insertContactData(d)
			if err != nil {
				http.Error(w, "Some Error Occureeeed", 5000)
				return
			}
			fmt.Printf("Data received and added in database\n", d)
			//data := "Hello Bhai sab badiya"
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(str)
		default:
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
		}
	} else {
		http.NotFound(w, r)
		return
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

func getUserData(id string) (User, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return User{}, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return User{}, err
	}
	collection := client.Database("test").Collection("test1")

	filter := bson.D{{"id", id}}

	var result User
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return User{}, err
	}

	return result, nil
}

func insertUserData(userData User) (interface{}, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return "", err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return "", err
	}
	collection := client.Database("test").Collection("test1")

	insertResult, err := collection.InsertOne(context.TODO(), userData)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID, nil
}

func insertContactData(contactData Contact) (interface{}, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return "", err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return "", err
	}
	collection := client.Database("test").Collection("test2")

	insertResult, err := collection.InsertOne(context.TODO(), contactData)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID, nil
}

func getContactData(userid string, cTime time.Time) ([]string, error) {

	var results []string

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return results, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return results, err
	}
	collection := client.Database("test").Collection("test2")

	filter := bson.D{{"$and", []bson.D{
		bson.D{{"contacttimestamp", bson.D{{"$and",
			[]bson.D{
				bson.D{{"$gt", cTime.Add(-14 * 24 * time.Hour)}},
				bson.D{{"$lt", cTime}},
			}}},
		}},
		bson.D{{"$or", []bson.D{
			bson.D{{"user1id", userid}},
			bson.D{{"user2id", userid}},
		}}},
	},
	}}

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return results, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Contact
		err = cur.Decode(&elem)
		if err != nil {
			return results, err
		}

		if elem.User1ID == userid {
			results = append(results, elem.User2ID)
		} else {
			results = append(results, elem.User1ID)
		}

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
	User1ID          string    `json:"user1id"`
	User2ID          string    `json:"user2id"`
	ContactTimestamp time.Time `json:"contacttimestamp"`
}
