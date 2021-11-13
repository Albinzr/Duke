package database

import (
	"context"
	util "duke/init/src/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoginDBConfig struct {
	Database       *mongo.Database
	CollectionName string
	Aud            string
	Iss            string
	collection     *mongo.Collection
	ctx            context.Context
}

type User struct {
	Username string `json" "username" bson: "username"`
	EmailId  string `json "emailId" bson: "emailId"`
	Password string `json "password" bson: "password"`
}

//Init :- initialize function
func (c *LoginDBConfig) Init() {
	c.collection = c.Database.Collection(c.CollectionName)
	indexName, err := c.collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		util.LogError("unable to create indexes for db", err)
		return
	}
	util.LogInfo("login module indexes created", indexName)

}

func (c *LoginDBConfig) CreateUser(user User) (primitive.ObjectID, error) {
	result, err := c.collection.InsertOne(c.ctx, user)
	if err != nil {
		util.LogError("unable insert value to db", err)
		return primitive.ObjectID{}, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	return id, nil
}

func (c *LoginDBConfig) FindUser(username string, password string) (primitive.M, error) {
	type Fields struct {
		username string `bson:"username"`
	}
	var result primitive.M
	filter := bson.M{"username": username}
	err := c.collection.FindOne(c.ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			util.LogError("no data found in db", err)
		} else if err != nil {
			util.LogError("unable to get data from db", err)
		}
		return nil, err
	}
	return result, nil
}
