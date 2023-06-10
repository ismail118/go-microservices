package data

import "go.mongodb.org/mongo-driver/mongo"

type Repository interface {
	Insert(entry LogEntry) error
	All() ([]*LogEntry, error)
	GetOne(id string) (*LogEntry, error)
	DropCollection(name string) error
	Update(l *LogEntry) (*mongo.UpdateResult, error)
}
