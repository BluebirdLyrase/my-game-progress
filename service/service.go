package service

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"my-game-progress/database"
	"my-game-progress/model/model_game"
)

func GetGameList() ([]model_game.Game, error) {
	gameCollection := database.DB.Collection("game")
	data, err := database.Fetch(gameCollection)
	if err != nil {
		log.Fatalf("Failed to fetch games: %v", err)
	}

	var games []model_game.Game
	for _, item := range data {
		var game model_game.Game
		// Decode bson.M into Game struct
		bsonBytes, err := bson.Marshal(item) // Convert bson.M to BSON bytes
		if err != nil {
			return nil, fmt.Errorf("failed to marshal bson.M: %w", err)
		}

		if err := bson.Unmarshal(bsonBytes, &game); err != nil {
			return nil, fmt.Errorf("failed to unmarshal to Game struct: %w", err)
		}
		games = append(games, game)
	}
	return games, nil

}
