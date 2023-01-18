package connectors

import (
	"context"
	"log"

	middlewares "orlangur.link/services/mini.note/handlers"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DbconnectMG -> connects mongo
func DbconnectMG() *mongo.Client {
	clientOptions := options.Client().ApplyURI(middlewares.DotEnvVariable("MONGO_URL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("⛒ Connection Failed to Database")
		log.Println(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("⛒ Connection Failed to Database")
		log.Println(err)
	} else {
		color.Green("⛁ Connected to Database Mongo")
	}
	return client
}
