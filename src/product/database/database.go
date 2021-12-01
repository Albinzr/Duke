package database

import (
	"context"
	util "duke/init/src/helpers"
	"duke/init/src/product/Config"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
)

type Config ProductConfig.Config

func (c *Config) Init() {
	c.Collection = c.Database.Collection(c.CollectionName)
	indexName, err := c.Collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "productId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		util.LogError("unable to create indexes for db", err)
		return
	}

	util.LogInfo("product module indexes created", indexName)
}

func (c *Config) Create(data url.Values) (primitive.ObjectID, error) {

	var value = make(map[string]interface{})

	for key, _ := range data {
		fmt.Println(key, data.Get(key))
		value[key] = data.Get(key)
	}
	print(value)
	productId, err := c.Collection.InsertOne(c.Ctx, value)

	if err != nil {
		util.LogError("unable to create productb", err)
		return primitive.ObjectID{}, err
	}
	id := productId.InsertedID.(primitive.ObjectID)
	return id, nil
}
