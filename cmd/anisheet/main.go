package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/Dophin2009/anisheet/pkg/data"
)

func main() {

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

	allMedia, _ := data.MediaGetAll(db)
	spew.Config = spew.ConfigState{
		Indent:                  "    ",
		DisableMethods:          true,
		DisableCapacities:       true,
		DisablePointerMethods:   true,
		DisablePointerAddresses: true,
	}
	spew.Config.Dump(allMedia)

	data.ClearDatabase(db)
}
