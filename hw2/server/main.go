package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"math/rand"
)

func main() {
	ctx := context.Background()
	mongoClient, err := createMongodbConnection(ctx, "mongodb://user:pass@localhost:27017")
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()

	// create new person
	app.Put("/person", func(ctx *fiber.Ctx) error {
		return nil
	})
	// return list of all persons
	app.Get("/person", func(ctx *fiber.Ctx) error {
		return nil
	})
	// get person by id
	app.Get("/person/:id", func(ctx *fiber.Ctx) error {
		return nil
	})
	// delete person by id
	app.Delete("/person/:id", func(ctx *fiber.Ctx) error {
		return nil
	})

	app.Get("/randomFile", func(ctx *fiber.Ctx) error {
		// first generate random string as file name
		// then create 1kb random file
		// get sha256 checksum of the file and sends it to client

		//err := ioutil.WriteFile("", nil, 0644)
		//if err != nil {
		//	return
		//}
		//getChecksumOfFile(path)
		return nil
	})

	err = app.Listen(":3000")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
func getChecksumOfFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	h.Write(content)
	sum := h.Sum(nil)
	return sum, nil
}

func generateRandomString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func createMongodbConnection(ctx context.Context, uri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, err
	}
	return client, nil
}
