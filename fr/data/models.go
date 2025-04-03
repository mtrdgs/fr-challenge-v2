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

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		QuoteEntry: QuoteEntry{},
	}
}

type Models struct {
	QuoteEntry QuoteEntry
}

type QuoteEntry struct {
	Carrier   []Carrier `bson:"carrier" json:"carrier"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// Carrier -
type Carrier struct {
	Name     string  `bson:"name" json:"name"`
	Service  string  `bson:"service" json:"service"`
	Deadline int     `bson:"deadline" json:"deadline"`
	Price    float64 `bson:"price" json:"price"`
}

func (q *QuoteEntry) Insert(entry QuoteEntry) error {
	collection := client.Database("fr").Collection("quotes")

	q.Carrier = []Carrier{}

	for _, value := range entry.Carrier {
		q.Carrier = append(q.Carrier, Carrier{
			Name:     value.Name,
			Service:  value.Service,
			Deadline: value.Deadline,
			Price:    value.Price,
		})
	}
	q.CreatedAt = time.Now()

	// save
	_, err := collection.InsertOne(context.TODO(), q)
	if err != nil {
		log.Println("Error inserting into quotes: ", err)
		return err
	}

	return nil
}

func (q *QuoteEntry) FindSpecific(amount int64) (quotes []QuoteEntry, err error) {
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
