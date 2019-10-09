package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/Dophin2009/anisheet/pkg/data"
)

func main() {
	spew.Config = spew.ConfigState{
		Indent:                  "    ",
		DisableMethods:          true,
		DisableCapacities:       true,
		DisablePointerMethods:   true,
		DisablePointerAddresses: true,
	}

	db, err := data.ConnectDatabase("/tmp/anisheet.db", true)
	if err != nil {
		panic("error connecting to database ")
	}
	defer db.Close()

	mitsuboshi := data.Media{
		Synopsis: "Three girls have fun",
		Titles: []data.Info{
			data.Info{
				Data:     "Mitsuboshi Colors",
				Language: "English",
			},
		},
	}
	err = data.MediaCreate(&mitsuboshi, db)
	if err != nil {
		fmt.Println(err)
	}

	notMitsuboshi := data.Media{
		Synopsis: "Three girls don't have fun",
	}
	err = data.MediaCreate(&notMitsuboshi, db)
	if err != nil {
		fmt.Println(err)
	}

	silverLink := data.Producer{
		Titles: []data.Info{
			data.Info{
				Data:     "Silver Link",
				Language: "English",
			},
		},
	}
	err = data.ProducerCreate(&silverLink, db)
	if err != nil {
		fmt.Println(err)
	}

	notSilverLink := data.Producer{
		Titles: []data.Info{
			data.Info{
				Data:     "Not Silver Link",
				Language: "English",
			},
		},
	}
	err = data.ProducerCreate(&notSilverLink, db)
	if err != nil {
		fmt.Println(err)
	}

	mitsuboshiSilverLink := data.MediaProducer{
		MediaID:    mitsuboshi.ID,
		ProducerID: silverLink.ID,
		Role:       "Studio",
	}
	err = data.MediaProducerCreate(&mitsuboshiSilverLink, db)
	if err != nil {
		fmt.Println(err)
	}

	notMitsuboshiNotSilverLink := data.MediaProducer{
		MediaID:    notMitsuboshi.ID,
		ProducerID: notSilverLink.ID,
		Role:       "Studio",
	}
	err = data.MediaProducerCreate(&notMitsuboshiNotSilverLink, db)
	if err != nil {
		fmt.Println(err)
	}

	mitsuboshiNotMitsuboshi := data.MediaRelation{
		OwnerID:      mitsuboshi.ID,
		RelatedID:    notMitsuboshi.ID,
		Relationship: "Side-story",
	}
	err = data.MediaRelationCreate(&mitsuboshiNotMitsuboshi, db)
	if err != nil {
		fmt.Println(err)
	}

	notMitsuboshiMitsuboshi := data.MediaRelation{
		OwnerID:      notMitsuboshi.ID,
		RelatedID:    mitsuboshi.ID,
		Relationship: "Main story",
	}
	err = data.MediaRelationCreate(&notMitsuboshiMitsuboshi, db)
	if err != nil {
		fmt.Println(err)
	}

	allMedia, err := data.MediaGetAll(db)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(allMedia)

	allProducers, err := data.ProducerGetAll(db)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(allProducers)

	allMediaProducers, err := data.MediaProducerGetAll(db)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(allMediaProducers)

	allMediaRelations, err := data.MediaRelationGetAll(db)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(allMediaRelations)

	data.ClearDatabase(db)
}
