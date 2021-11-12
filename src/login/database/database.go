package database

import (
	"context"
	util "duke/init/src/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoginDBConfig struct {
	URL            string
	DatabaseName   string
	CollectionName string
	//
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

type User struct {
	Username string `json" "username" bson: "username"`
	EmailId  string `json "emailId" bson: "emailId"`
	Password string `json "password" bson: "password"`
}

//Init :- initalize function
func (c *LoginDBConfig) Init() {
	var err error
	c.client, err = mongo.NewClient(options.Client().ApplyURI(c.URL))

	util.LogError("databaseClientError", err)
	c.database = c.client.Database(c.DatabaseName)
	c.collection = c.database.Collection(c.CollectionName)
	//var cancel context.CancelFunc
	c.ctx, _ = context.WithCancel(context.Background())
	err = c.client.Connect(c.ctx)
	util.LogError("databaseConnectionError", err)

	if err != nil {
		util.LogError("Database connection issue", err)
		return
	} else {
		util.LogInfo("Database connected")
	}
	//defer cancel()
}

func (c *LoginDBConfig) CreateUser(user User) error {
	result, err := c.collection.InsertOne(c.ctx, user)
	util.LogInfo(result.InsertedID, err)
	return err
}
