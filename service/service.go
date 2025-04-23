package service

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"my-game-progress/database"
	model_environment "my-game-progress/model/environment"
	"my-game-progress/model/model_game"
)

func GetGameList(filter bson.M, sort bson.M, limit int64) ([]model_game.GameBase, error) {
	gameCollection := database.DB.Collection("game")
	opts := options.Find().SetSort(sort).SetLimit(limit)
	cursor, err := gameCollection.Find(database.Context, filter, opts)
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

func GetGameFullDetail(filter bson.M, sort bson.M, limit int64) ([]model_game.Game, error) {
	var games []model_game.Game
	gameCollection := database.DB.Collection("game")
	opts := options.Find().SetSort(sort).SetLimit(limit)
	cursor, err := gameCollection.Find(database.Context, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch games: %w", err)
	}
	defer cursor.Close(database.Context)

	if err := cursor.All(database.Context, &games); err != nil {
		return nil, fmt.Errorf("failed to decode games: %w", err)
	}

	return games, nil
}

func GetEnvironmentsList() ([]model_environment.Environment, error) {
	var environments []model_environment.Environment
	collection := database.DB.Collection("environment")

	opts := options.Find().SetProjection(bson.M{"pc_spec": 0})

	cursor, err := collection.Find(database.Context, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query environments: %w", err)
	}
	defer cursor.Close(database.Context)

	for cursor.Next(database.Context) {
		var env model_environment.Environment
		if err := cursor.Decode(&env); err != nil {
			return nil, fmt.Errorf("failed to decode environment: %w", err)
		}
		environments = append(environments, env)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return environments, nil
}
