package service

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"my-game-progress/database"
	"my-game-progress/model/model_game"
)

func GetGameList() ([]model_game.GameBase, error) {
	gameCollection := database.DB.Collection("game")
	cursor, err := gameCollection.Find(database.Context, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch games: %w", err)
	}
	defer cursor.Close(database.Context)

	var games []model_game.GameBase
	if err := cursor.All(database.Context, &games); err != nil {
		return nil, fmt.Errorf("failed to decode games: %w", err)
	}

	return games, nil
}
