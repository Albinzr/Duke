package database

import (
	"context"
	util "duke/init/src/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Config for database connection
type Config struct {
	URL          string
	DatabaseName string

	client   *mongo.Client
	database *mongo.Database
	ctx      context.Context
}

//Init :- initalize function
func (c *Config) Init() {
	var err error
	c.client, err = mongo.NewClient(options.Client().ApplyURI(c.URL))

	util.LogError("databaseClientError", err)
	c.database = c.client.Database(c.DatabaseName)

	c.ctx, _ = context.WithCancel(context.Background())
	err = c.client.Connect(c.ctx)
	util.LogError("databaseConnectionError", err)

	if err != nil {
		util.LogError("Database connection issue", err)
		return
	} else {
		util.LogInfo("Database connected")
	}

}
