package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var client *mongo.Client

type MongoRepository struct {
	Conn *mongo.Client
}

func NewMongoRepository(mongo *mongo.Client) *MongoRepository {
	client = mongo
	return &MongoRepository{
		Conn: mongo,
	}
}

// func New(mongo *mongo.Client) Models {
// 	client = mongo

// 	return Models{
// 		QuoteEntry: QuoteEntry{},
// 	}
// }

// type Models struct {
// 	QuoteEntry QuoteEntry
// }

// QuoteEntry -
type QuoteEntry struct {
	Carrier   []Carrier  `bson:"carrier" json:"carrier"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at,omitempty"`
}

// Carrier -
type Carrier struct {
	Name     string  `bson:"name" json:"name"`
	Service  string  `bson:"service" json:"service"`
	Deadline int     `bson:"deadline" json:"deadline"`
	Price    float64 `bson:"price" json:"price"`
}

// Insert - stores quotes from freterapido api in bd
func (q *MongoRepository) Insert(entry QuoteEntry) error {
	collection := client.Database("fr").Collection("quotes")

	currentTime := time.Now()
	entry.CreatedAt = &currentTime

	// save
	_, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		log.Println("Error inserting into quotes: ", err)
		return err
	}

	return nil
}

// FindSpecific - gets list of quotes from db
func (q *MongoRepository) FindSpecific(amount int64) (quotes []QuoteEntry, err error) {
	var cursor *mongo.Cursor

	collection := client.Database("fr").Collection("quotes")

	if amount > 0 {
		cursor, err = collection.Find(context.TODO(), bson.M{}, options.Find().SetLimit(amount).SetSort(bson.M{"created_at": -1}))
	} else {
		cursor, err = collection.Find(context.TODO(), bson.M{})
	}

	if err != nil {
		log.Printf("Error retrieving quotes: %v", err)
		return quotes, err
	}

	// convert cursor into array
	err = cursor.All(context.TODO(), &quotes)
	if err != nil {
		log.Printf("Error converting quotes into JSON: %v", err)
		return quotes, err
	}

	return quotes, err
}
