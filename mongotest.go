package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID                string `json:"string"`
	Name              string `json:"name"`
	Dob               string `json:"string"`
	PhoneNo           int    `json:"int"`
	Email             string `json:"string"`
	CreationTimestamp time.Time
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("test1")

	/*user1 := User{"1", "Ash", "2000-10-02", 7849322792, "user1@gmail.com", time.Now()}
	user2 := User{"2", "Brock", "1998-03-02", 7843252792, "user2@gmail.com", time.Now()}
	user3 := User{"3", "Mistry", "1998-03-01", 7843255243, "user3@gmail.com", time.Now()}

	insertResult, err := collection.InsertOne(context.TODO(), user1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	trainers := []interface{}{user2, user3}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)*/

	/*filter := bson.D{{"id", 2}}

	var result User

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)*/

	// Here's an array in which you can store the decoded documents

	findTime := time.Now()
	findTime.Add(-2 * time.Hour)

	var results []string

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{"creationtimestamp", bson.D{{"$lt", findTime}}}})
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem.ID)
	}

	fmt.Printf("Found a single document: %+v\n", results)

	/*if err = collection.Drop(ctx); err != nil {
		log.Fatal(err)
	}*/

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
