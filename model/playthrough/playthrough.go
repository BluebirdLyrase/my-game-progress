package playthrough

import "go.mongodb.org/mongo-driver/bson/primitive"

type Playthrough struct {
	ID            *primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Difficulty    string              `json:"difficulty"`
	DateFinished  string              `json:"dateFinished"` //* DD/MM/YYYY
	Remark        string              `json:"remark"`
	Screenshots   []string            `json:"screenshots"`
	GameID        *primitive.ObjectID `json:"game_id" bson:"_id,omitempty"`
	EnvironmentID *primitive.ObjectID `json:"environment_id" bson:"_id,omitempty"`
}

type PlaythroughInputParam struct {
	Difficulty    string              `json:"difficulty"`
	DateFinished  string              `json:"dateFinished"` //* DD/MM/YYYY
	Remark        string              `json:"remark"`
	GameID        *primitive.ObjectID `json:"game_id" bson:"_id,omitempty"`
	EnvironmentID *primitive.ObjectID `json:"environment_id" bson:"_id,omitempty"`
}
