package main

import (
	"fmt"

	"github.com/wangkekekexili/gops/model"
	"github.com/wangkekekexili/gops/util"
	"gopkg.in/mgo.v2"
)

func main() {
	// Connect to mlab mongodb.
	uri, err := util.BuildMongodbURI()
	if err != nil {
		panic(err)
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("gops").C("gops")
	gamestop := &gops.Gamestop{}
	gamesToInsert, gamesToUpdate, err := gamestop.GetGamesInfo(c)
	if err != nil {
		panic(err)
	}
	if len(gamesToInsert) > 0 {
		fmt.Printf("inserting %d documents\n", len(gamesToInsert))
		if err = c.Insert(gamesToInsert...); err != nil {
			panic(err)
		}
	}
	if len(gamesToUpdate) > 0 {
		fmt.Printf("updating %d documents\n", len(gamesToUpdate))
		for objectID, game := range gamesToUpdate {
			if err = c.UpdateId(objectID, game); err != nil {
				fmt.Println(err)
			}
		}
	}
}
