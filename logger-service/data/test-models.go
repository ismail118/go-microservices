package data

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTestRepository struct {
	Client *mongo.Client
}

func NewMongoTestRepository(c *mongo.Client) *MongoTestRepository {
	return &MongoTestRepository{
		Client: c,
	}
}

func (r *MongoTestRepository) Insert(entry LogEntry) error {
	return nil
}

func (r *MongoTestRepository) All() ([]*LogEntry, error) {
	var logs []*LogEntry
	return logs, nil
}

func (r *MongoTestRepository) GetOne(id string) (*LogEntry, error) {
	var entry LogEntry
	return &entry, nil
}

func (r *MongoTestRepository) DropCollection(name string) error {
	return nil
}

func (r *MongoTestRepository) Update(l *LogEntry) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult
	return result, nil
}
