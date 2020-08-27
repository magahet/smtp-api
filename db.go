package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Timeout operations after N seconds
	connectTimeout           = 5
	connectionStringTemplate = "mongodb://%s:%s@%s/penpal"
)

// GetConnection - Retrieves a client to the DocumentDB
func getConnection() (*mongo.Client, context.Context, context.CancelFunc, error) {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	clusterEndpoint := os.Getenv("DB_SERVER")

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
	} else {
		fmt.Println("Connected to MongoDB!")
	}

	return client, ctx, cancel, err
}

// GetAllEvents Retrives all events from the db
func GetAllEvents() ([]*Event, error) {
	var events []*Event

	client, ctx, cancel, err := getConnection()
	if err != nil {
		return events, err
	}
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("penpal")
	collection := db.Collection("events")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &events)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return events, nil
}

// GetEventByID Retrives a event by its id from the db
func GetEventByID(id primitive.ObjectID) (*Event, error) {
	var event *Event

	client, ctx, cancel, err := getConnection()
	if err != nil {
		return event, err
	}
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("penpal")
	collection := db.Collection("events")
	result := collection.FindOne(ctx, bson.D{})
	if result == nil {
		return nil, errors.New("Could not find a Event")
	}
	err = result.Decode(&event)

	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	log.Printf("Events: %v", event)
	return event, nil
}

//Create creating a event in a mongo
func Create(event *Event) (primitive.ObjectID, error) {
	client, ctx, cancel, err := getConnection()
	if err != nil {
		return primitive.NilObjectID, err
	}
	defer cancel()
	defer client.Disconnect(ctx)
	event.ID = primitive.NewObjectID()

	result, err := client.Database("penpal").Collection("events").InsertOne(ctx, event)
	if err != nil {
		log.Printf("Could not create Event: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

//Update updating an existing event in a mongo
func Update(event *Event) (*Event, error) {
	var updatedEvent *Event
	if event.ID == primitive.NilObjectID {
		return updatedEvent, errors.New("Event ID not provided")
	}
	client, ctx, cancel, err := getConnection()
	if err != nil {
		return updatedEvent, err
	}
	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": event,
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err = client.Database("penpal").Collection("events").FindOneAndUpdate(ctx, bson.M{"_id": event.ID}, update, &opt).Decode(&updatedEvent)
	if err != nil {
		log.Printf("Could not save Event: %v", err)
		return nil, err
	}
	return updatedEvent, nil
}
