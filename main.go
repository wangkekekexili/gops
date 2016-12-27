package main

import (
	"github.com/uber-go/zap"
	"github.com/wangkekekexili/gops/model"
	"github.com/wangkekekexili/gops/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	// Config logger.
	logger := zap.New(zap.NewJSONEncoder(zap.NoTime()))

	// Connect to mlab mongodb.
	uri, err := util.BuildMongodbURI()
	if err != nil {
		logger.Error("Cannot build MongoDB URI.", zap.String("err", err.Error()))
		util.SendAlert()
		return
	}
	logger.Info("Connecting to MongoDB.", zap.String("uri", uri))
	session, err := mgo.Dial(uri)
	if err != nil {
		logger.Error("Cannot connect to MongoDB.", zap.String("err", err.Error()))
		util.SendAlert()
		return
	}
	defer session.Close()

	c := session.DB("gops").C("gops")
	for _, handler := range gops.GetAllGameFetchers() {
		logger.Info("Handler starts.",
			zap.String("source", handler.GetSource()),
		)
		gameNames, gamesByNameAndCondition, err := handler.GetGamesInfo(c)
		if err != nil {
			logger.Error("Unable to get games info.",
				zap.String("source", gops.ProductSourceGamestop),
				zap.String("err", err.Error()),
			)
			util.SendAlert()
			continue
		}

		// Get existing documents and see if we need to update them.
		namesSubQuery := make([]bson.M, len(gameNames))
		for i, gameName := range gameNames {
			namesSubQuery[i] = bson.M{"name": gameName}
		}
		cursor := c.Find(bson.M{"source": handler.GetSource(), "$or": namesSubQuery}).Iter()
		var result gops.BasicGameInfo
		for cursor.Next(&result) {
			name := result.Name
			condition := result.Condition
			key := name + condition
			if gameNewPricePoint, ok := gamesByNameAndCondition[key]; ok {
				// Skip if the price is not changed.
				if result.PriceHistory[len(result.PriceHistory)-1].Price == gameNewPricePoint.GetRecentPrice() {
					delete(gamesByNameAndCondition, key)
					continue
				}
				// Make a copy of result.
				gameToUpdate := result
				gameToUpdate.PriceHistory = append(gameToUpdate.PriceHistory, gameNewPricePoint.GetPriceHistory()[0])
				gamesByNameAndCondition[key] = &gameToUpdate
			}
		}
		cursor.Close()

		var gamesToInsert []interface{}
		gamesToUpdate := make(map[bson.ObjectId]interface{})

		for _, game := range gamesByNameAndCondition {
			if game.GetID().Hex() == "" {
				gamesToInsert = append(gamesToInsert, game)
			} else {
				gamesToUpdate[game.GetID()] = game
			}
		}

		if len(gamesToInsert) > 0 {
			logger.Info("Inserting documents.", zap.Int("number", len(gamesToInsert)))
			if err = c.Insert(gamesToInsert...); err != nil {
				logger.Warn("Unable to insert games.",
					zap.String("err", err.Error()),
				)
				util.SendAlert()
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
				util.SendAlert()
			}
		}
		logger.Info("Handler finishes.",
			zap.String("source", handler.GetSource()),
		)
	}
}
