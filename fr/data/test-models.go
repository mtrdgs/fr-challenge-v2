package data

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTestRepository struct {
	Conn *mongo.Client
}

func NewMongoTestRepository(mongo *mongo.Client) *MongoTestRepository {
	return &MongoTestRepository{
		Conn: mongo,
	}
}

func (q *MongoTestRepository) Insert(entry QuoteEntry) error {
	return nil
}

func (q *MongoTestRepository) FindSpecific(amount int64) (quotes []QuoteEntry, err error) {
	currentTime := time.Now()
	quotes = []QuoteEntry{
		{
			Carrier: []Carrier{
				{
					Name:     "test",
					Service:  "test",
					Deadline: 1,
					Price:    1.5,
				},
			},
			CreatedAt: &currentTime,
		},
	}

	return quotes, err
}
