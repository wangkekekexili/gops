package main

import (
	"github.com/uber-go/zap"
	"github.com/wangkekekexili/gops/model"
	"github.com/wangkekekexili/gops/util"
	"gopkg.in/mgo.v2"
)

func main() {
	// Config logger.
	logger := zap.New(zap.NewJSONEncoder(zap.NoTime()))

	// Connect to mlab mongodb.
	uri, err := util.BuildMongodbURI()
	if err != nil {
		logger.Fatal("Cannot build MongoDB URI.", zap.String("message", err.Error()))
		return
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		logger.Fatal("Cannot connect to MongoDB.", zap.String("message", err.Error()))
	}
	defer session.Close()

	c := session.DB("gops").C("gops")
	gamestop := &gops.Gamestop{}
	gamesToInsert, gamesToUpdate, err := gamestop.GetGamesInfo(c)
	if err != nil {
		logger.Error("Unable to get games info.", zap.String("source", gops.ProductSourceGamestop))
	}
	if len(gamesToInsert) > 0 {
		logger.Info("Inserting documents.", zap.Int("number", len(gamesToInsert)))
		if err = c.Insert(gamesToInsert...); err != nil {
			logger.Warn("Unable to insert games.")
		}
	}
	if len(gamesToUpdate) > 0 {
		logger.Info("Updating documents.", zap.Int("number", len(gamesToUpdate)))
		for objectID, game := range gamesToUpdate {
			if err = c.UpdateId(objectID, game); err != nil {
				logger.Warn("Unable to update a document.",
					zap.String("id", objectID.Hex()),
					zap.Object("game", game),
				)
			}
		}
	}
}