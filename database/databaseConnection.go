package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)


func DBinstance() *mongo.Client{

	mongoDbUri := "mongodb+srv://karthik:nair12345678@cluster0.osthylz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	opts := options.Client().ApplyURI(mongoDbUri);

	client, err := mongo.Connect(opts);
	if err != nil {
		log.Fatal("error while connecting to db : " , err)
	  }
	  defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
		  log.Fatal("error happened while connecting to database : ", err);
		}
	  }()
	  // Send a ping to confirm a successful connection
	  if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("error while pinging db: ", err)
	  }
	  fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	  return client
	
}


var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)
	return collection
}
  