package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/gofiber/fiber/v2/middleware/logger"
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
	app.Use(logger.New())

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
		fileName := fmt.Sprintf("%s.txt", string(generateRandomString(10)))
		path := fmt.Sprintf("/serverdata/%s.txt", fileName)
		content := generateRandomString(1024)
		err = ioutil.WriteFile(path, content, 0644)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			err = ctx.JSON(fiber.Map{"error": "can not serve file"})
			return err
		}
		file, err := getChecksumOfFile(path)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not serve file"})
			return err
		}
		ctx.Set("checksum", file)
		ctx.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		err = ctx.SendFile(path)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			ctx.JSON(fiber.Map{"error": "can not serve file"})
			return err
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
