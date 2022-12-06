package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"example/typed/pkg/domain"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("typed").Collection("messages")

	// now let's save a message to the collection
	msg := domain.Email{
		Title:  "Dear Correspondant,",
		Conent: "This email is not intended for you.",
	}

	snd_gob, err := msg.GetGob()
	if err != nil {
		panic(err)
	}

	_, err = coll.InsertOne(context.TODO(), snd_gob)
	if err != nil {
		panic(err)
	}

	// retrieve the document
	var ret_gob domain.MessageGob
	err = coll.FindOne(context.TODO(), bson.M{"_id": snd_gob.ID}).Decode(&ret_gob)
	if err != nil {
		panic(err)
	}

	ret_msg, err := ret_gob.GetMessage()
	if err != nil {
		panic(err)
	}

	content, err := ret_msg.Marshal()
	if err != nil {
		panic(err)
	}

	fmt.Print(string(content))

}
