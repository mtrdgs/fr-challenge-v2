package data

import (
	"context"
	"log"
	"time"
)

type quoteEntry struct {
	Carrier   []Carrier `json:"carrier"`
	CreatedAt time.Time `json:"created_at"`
}

// Carrier -
type Carrier struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline int     `json:"deadline"`
	Price    float64 `json:"price"`
}

func InsertDB(entry quoteEntry) (insert quoteEntry, err error) {
	collection := client.Database("fr").Collection("quotes")

	for _, value := range entry.Carrier {
		insert.Carrier = append(insert.Carrier, Carrier{
			Name:     value.Name,
			Service:  value.Service,
			Deadline: value.Deadline,
			Price:    value.Price,
		})
	}
	insert.CreatedAt = time.Now()

	// save
	_, err = collection.InsertOne(context.TODO(), insert)
	if err != nil {
		log.Println("Error inserting into quotes: ", err)
		return insert, err
	}

	return insert, nil
}
