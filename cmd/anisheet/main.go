package main

import (
	"encoding/json"
	"fmt"

	"gitlab.com/Dophin2009/anisheet/pkg/data"
)

func main() {
	database := data.ConnectWithMigrations("/tmp/anisheet.db")

	mitsuboshi := data.Media{
		Synopsis: "The show follows the fun activities of three grade-school girls.",
		Titles: []data.Title{
			data.Title{Name: "Mitsuboshi Colors", Language: "English"},
		},
	}
	data.MediaCreate(&mitsuboshi, database)

	notMitsuboshi := data.Media{
		Synopsis: "The show doesn't follow the fun activities of three grade-school girls.",
		Titles: []data.Title{
			data.Title{Name: "Not Mitsuboshi Colors", Language: "English"},
		},
	}
	data.MediaCreate(&notMitsuboshi, database)

	relation := data.MediaRelation{
		Owner:    mitsuboshi.ID,
		Related:  notMitsuboshi.ID,
		Relation: "Side-story",
	}
	err := data.MediaRelationCreate(&relation, database)
	if err != nil {
		panic(err)
	}

	silverLink := data.Producer{
		Titles: []data.Title{
			data.Title{Name: "Silver Link", Language: "English"},
		},
	}
	err = data.ProducerCreate(&silverLink, database)
	if err != nil {
		panic(err)
	}

	producerRelation := data.MediaProducer{
		MediaID:    mitsuboshi.ID,
		ProducerID: silverLink.ID,
		Role:       "Studio",
	}
	err = data.MediaProducerCreate(&producerRelation, database)
	if err != nil {
		panic(err)
	}

	found, err := data.ProducerGetByID(1, database)
	if err != nil {
		panic(err)
	}
	s, _ := json.MarshalIndent(found, "", "\t")
	fmt.Print(string(s))

}
