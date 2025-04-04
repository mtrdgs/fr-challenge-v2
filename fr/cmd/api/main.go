package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mtrdgs/fr/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultPort = "80"
	mongoURL    = "mongodb://mongo:27017"
	apiURL      = "https://sp.freterapido.com/api/v3/quote/simulate"
)

var client *mongo.Client

type Config struct {
	Repo   data.RepositoryPattern
	Client *http.Client
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create ctx in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Client: &http.Client{},
		//Models: data.New(client),
	}

	app.setUpRepo(client)

	log.Printf("Starting server on port %s.", defaultPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", defaultPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create client connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panicln("Error connecting to mongo: ", err)
		return nil, err
	}
	log.Println("Connected to Mongo!")

	return c, nil
}

func (app *Config) setUpRepo(conn *mongo.Client) {
	mongo := data.NewMongoRepository(conn)
	app.Repo = mongo
}
