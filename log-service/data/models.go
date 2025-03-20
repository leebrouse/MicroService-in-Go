package data

// using mongo db
import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Struct
type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID       string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string    `bson:"name" json:"name"`
	Data     string    `bson:"data" json:"data"`
	CreateAt time.Time `bson:"create_at" json:"create_at"`
	UpdateAt time.Time `bson:"update_at" json:"update_at"`
}

// Function
func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

// Insert one LogEntry
func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:     entry.Name,
		Data:     entry.Data,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs", err)
		return err
	}
	return nil
}

// Show all of the log in the mongo
func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "create_at", Value: -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}

	}

	return logs, nil
}

// Get one specific log by id
func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	// set context for limiting the operation time
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//get log collection
	collection := client.Database("logs").Collection("logs")

	//get id
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	// set context for limiting the operation time
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//get log collection
	collection := client.Database("logs").Collection("logs")

	//drop the collection
	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

// Update one specific log by id
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//get log collection
	collection := client.Database("logs").Collection("logs")

	//get id
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	//call update function
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: l.Name},
				{Key: "data", Value: l.Data},
				{Key: "update_at", Value: time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
