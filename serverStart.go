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

	if strings.Contains(r.URL.Path, "/users") {

		switch r.Method {
		case "GET":
			id := strings.SplitN(r.URL.Path, "/", 3)[2]
			d, err := getUserData(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break
			}
			fmt.Printf("data for id %d send\n", id)

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(d)
		case "POST":
			var d User

			if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			str, err := insertUserData(d)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Printf("\n", d)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(str)
		default:
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
		}
	} else if strings.Contains(r.URL.Path, "/contacts") {
		fmt.Println("\nyaha tak 1")

		switch r.Method {
		case "GET":

			var userid string
			var currentTime time.Time
			for k, v := range r.URL.Query() {
				fmt.Println("\nyaha tak 2")
				//ct := time.Now()
				if k == "user" {
					userid = v[0]
				} else if k == "infection_timestamp" {
					fmt.Printf("%T\n", v[0])
					currentTime, _ = time.Parse("2014-11-12T11:45:26.371Z", v[0])
				} else {
					http.Error(w, "Incorrect Naming Convention ", http.StatusInternalServerError)
					return
				}
			}
			fmt.Println("\nyaha tak 3 %s", userid)
			fmt.Println(currentTime)
			//contactTime, _ := time.Parse("2020-10-2", ct)
			victims, err := getContactData(userid, currentTime)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("Possible Covid victims can be %+v\n", victims)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(victims)

		case "POST":
			var d Contact

			if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
				http.Error(w, "error1", http.StatusInternalServerError)
				return
			}
			str, err := insertContactData(d)
			if err != nil {
				http.Error(w, "error2", http.StatusInternalServerError)
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
		http.Error(w, "Nahi Mila", 404)
		return
	}

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
	fmt.Println("\nyaha tak 4", cTime)
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

	//cTime2 := cTime.Add(-14 * 24 * time.Hour)
	kTime := time.Now().Add(-14 * 24 * time.Hour)
	fmt.Println(kTime)
	kTime2 := time.Now()
	//filter := bson.D{{"$and", []bson.D{
	/*filter := bson.D{{"contacttimestamp", bson.D{{"$and",
		[]bson.D{
			bson.D{{"$gt", kTime}},
			bson.D{{"$lt", kTime2}},
		}}},
	}}*/
	/*bson.D{{"$or", []bson.D{
			bson.D{{"user1id", userid}},
			bson.D{{"user2id", userid}},
		}}},
	},
	}}*/
	filter := bson.D{{"$and", []bson.D{bson.D{{"contacttimestamp", bson.D{{"$gt", kTime}, {"$lt", kTime2}}}},
		bson.D{{"$or", []bson.D{bson.D{{"user1id", userid}}, bson.D{{"user2id", userid}}}}}}}}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return results, err
	}
	fmt.Println("\nyaha tak 5")

	for cur.Next(context.TODO()) {

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
	http.HandleFunc("/", receive)
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
