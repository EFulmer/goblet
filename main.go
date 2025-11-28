package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	Name string `json:"name"`
}

type WeighIn struct {
	User      string    `json:"user,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Weight    *Weight   `json:"weight,omitempty"`
}

type Weight struct {
	Quantity float64       `json:"quantity"`
	Unit     UnitOfMeasure `json:"string"`
}

type UnitOfMeasure string

type Workout struct {
	User      string      `json:"user,omitempty"`
	Date      string      `json:date` // TODO understand the Golang Time type in order to deserialize this into just the date.
	Exercises []*Exercise `json:"exercises"`
}

type Exercise struct {
	Name string `json:"name"`
	Sets []*Set `json:"sets"`
}

type Set struct {
	Weight *Weight `json:"weight"`
	Reps   uint    `json:"reps"`
}

const (
	Pounds    UnitOfMeasure = "pounds"
	Kilograms UnitOfMeasure = "kilograms"
)

func main() {
	defer os.Exit(0)
	var demo bool
	flag.BoolVar(&demo, "demo", false, "Print out a command line demo.")
	flag.Parse()
	fmt.Println("Hello, world!")
	if demo {
		addTestDataToMongo("test_db")
	}
}

func readWeighInsMongo() {
	uri := "mongodb://root:example@localhost:27017/"
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Println("error! %q", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("test_db").Collection("weigh_ins")

	var result bson.M

	err = coll.FindOne(context.TODO(), bson.D{{}}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		// fmt.Printf("No document was found with the title %s\n", title)
		fmt.Printf("No document was found\n")
		return
	}
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}

func addTestDataToMongo(dbName string) {
	rawUsers, err := os.ReadFile("data/users.json")
	if err != nil {
		log.Fatalf("Error reading data/users.json: %#v", err)
	}
	var users []User
	err = json.Unmarshal(rawUsers, &users)
	if err != nil {
		log.Fatalf("Error unmarshaling data/users.json: %#v", err)
	}
	fmt.Printf("%d users found.\n", len(users))
	for _, user := range users {
		fmt.Printf("\t%#v\n", user)
	}
	fmt.Println("")

	rawWeighIns, err := os.ReadFile("data/weighins.json")
	if err != nil {
		log.Fatalf("Error reading data/weighins.json: %#v", err)
	}
	var weighIns []WeighIn
	err = json.Unmarshal(rawWeighIns, &weighIns)
	if err != nil {
		log.Fatalf("Error unmarshaling data/weighins.json: %#v", err)
	}
	fmt.Printf("%d weigh-ins found.\n", len(weighIns))
	for _, weighIn := range weighIns {
		fmt.Printf("\t%#v\n", weighIn)
	}
	fmt.Println("")

	rawWorkouts, err := os.ReadFile("data/workouts.json")
	if err != nil {
		log.Fatalf("Error reading data/workouts.json: %#v", err)
	}
	var workouts []Workout
	err = json.Unmarshal(rawWorkouts, &workouts)
	if err != nil {
		log.Fatalf("Error unmarshaling data/workouts.json: %#v", err)
	}
	fmt.Printf("%d workouts found.\n", len(workouts))
	for _, workout := range workouts {
		fmt.Printf("\t%#v\n", workout)
	}
	fmt.Println("")

	fmt.Printf("Inserting all documents into MongoDB DB %s\n", dbName)
	fmt.Println("Inserting users...")
	addDocuments(dbName, "users", "user", users)
	fmt.Println("Inserting weigh-ins...")
	addDocuments(dbName, "weigh_ins", "weigh-in", weighIns)
	fmt.Println("Inserting workouts...")
	addDocuments(dbName, "workouts", "workout", workouts)
}

func addDocuments[T any](dbName, collectionName, documentTypeName string, docs []T) {
	// TODO: Discover and apply the idiomatic Golang method for reading sensitive data from .env files here.
	uri := "mongodb://root:example@localhost:27017/"
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("error! %q", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database(dbName).Collection(collectionName)

	result, err := coll.InsertMany(context.TODO(), docs)
	fmt.Printf("Number of %s inserted: %v\n", documentTypeName, len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		fmt.Printf("Inserted %s with _id: %v\n", documentTypeName, id)
	}
}
