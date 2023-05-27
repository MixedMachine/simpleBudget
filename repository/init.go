package repository

import (
	"context"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitializeDB() (*context.Context, *mongo.Client) {
	ctx := context.TODO()

	opts := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("mongo connection established")

	return &ctx, client
}
