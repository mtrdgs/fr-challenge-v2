package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
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
