package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"

	"io/ioutil"
	"math/rand"
)

const USER_DATABASE = "user"
const PERSON_COLLECTION = "person"

type Person struct {
	Name   string `json:"name"`
	Family string `json:"family"`
}

func main() {
	mongoUsername, ok := os.LookupEnv("MONGO_USERNAME")
	if !ok {
		panic("set MONGO_USERNAME ")
	}
	mongoPassword, ok := os.LookupEnv("MONGO_PASSWORD")
	if !ok {
		panic("set MONGO_PASSWORD ")
	}
	mongoHost, ok := os.LookupEnv("MONGO_HOST")
	if !ok {
		panic("set MONGO_HOST ")
	}
	mongoPort, ok := os.LookupEnv("MONGO_PORT")
	if !ok {
		panic("set MONGO_PORT ")
	}

	ctx := context.Background()
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUsername, mongoPassword, mongoHost, mongoPort)
	mongoClient, err := createMongodbConnection(ctx, url)
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()
	app.Use(logger.New())

	// create new person
	app.Put("/person", func(ctx *fiber.Ctx) error {
		person := Person{}
		err := ctx.BodyParser(&person)

		if err != nil {
			ctx.Status(fiber.StatusBadRequest)
			ctx.JSON(fiber.Map{"error": "bad request"})
			return nil
		}
		fmt.Println("person", person)
		res, err := mongoClient.Database(USER_DATABASE).Collection(PERSON_COLLECTION).InsertOne(context.Background(), person)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not create person"})
		}
		ctx.JSON(fiber.Map{"id": res.InsertedID})
		return nil
	})
	// return list of all persons
	app.Get("/person/all", func(ctx *fiber.Ctx) error {
		var persons []bson.M
		res, err := mongoClient.Database(USER_DATABASE).Collection(PERSON_COLLECTION).Find(context.Background(), bson.M{})
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not get persons"})
			return nil
		}
		err = res.All(context.Background(), &persons)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not get persons"})
			return nil
		}
		fmt.Println("persons: ", persons)
		ctx.JSON(persons)
		return nil
	})
	// get person by id
	app.Get("/person/:id", func(ctx *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(utils.CopyString(ctx.Params("id")))
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not get persons"})
			return nil
		}
		var person Person
		err = mongoClient.Database(USER_DATABASE).Collection(PERSON_COLLECTION).FindOne(context.Background(), bson.M{"_id": id}).Decode(&person)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not get person"})
			return nil
		}
		ctx.JSON(person)
		return nil
	})
	// delete person by id
	app.Delete("/person/:id", func(ctx *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(utils.CopyString(ctx.Params("id")))
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not delete person"})
			return nil
		}
		var person Person
		err = mongoClient.Database(USER_DATABASE).Collection(PERSON_COLLECTION).FindOneAndDelete(context.Background(), bson.M{"_id": id}).Decode(&person)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not get persons"})
			return nil
		}
		ctx.JSON(person)
		return nil
	})

	app.Get("/randomFile", func(ctx *fiber.Ctx) error {
		fileName := fmt.Sprintf("%s.txt", string(generateRandomString(10)))
		path := fmt.Sprintf("/serverdata/%s.txt", fileName)
		content := generateRandomString(1024)
		err = ioutil.WriteFile(path, content, 0644)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			err = ctx.JSON(fiber.Map{"error": "can not serve file"})
			return nil
		}
		file, err := getChecksumOfFile(path)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not serve file"})
			return nil
		}
		ctx.Set("checksum", file)
		ctx.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		err = ctx.SendFile(path)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not serve file"})
			return nil
		}
		return nil
	})

	err = app.Listen(":3000")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
func getChecksumOfFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write(content)
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}

func generateRandomString(n int) []byte {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
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
