package database

import (
	util "duke/init/src/helpers"
	"duke/init/src/product/Config"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
)

type Config ProductConfig.Config

func (c *Config) Init() {
	c.Collection = c.Database.Collection(c.CollectionName)
	//indexName, err := c.Collection.Indexes().CreateOne(
	//	context.Background(),
	//	mongo.IndexModel{
	//		Keys:    bson.D{{Key: "productId", Value: 1}},
	//		Options: options.Index().SetUnique(true),
	//	},
	//)
	//if err != nil {
	//	util.LogError("unable to create indexes for db", err)
	//	return
	//}
	//
	//util.LogInfo("product module indexes created", indexName)
}

func (c *Config) Create(data url.Values) (primitive.ObjectID, error) {

	var value = make(map[string]interface{})
	for key, _ := range data {
		value[key] = data.Get(key)
	}
	productId, err := c.Collection.InsertOne(c.Ctx, value)
	if err != nil {
		util.LogError("unable to create product", err)
		return primitive.ObjectID{}, err
	}
	id := productId.InsertedID.(primitive.ObjectID)
	return id, nil
}
func (c *Config) Update(data url.Values) (bool, error) {

	var value = make(map[string]interface{})
	var filter bson.M
	for key, _ := range data {
		if key == "_id"{
			_id, err := primitive.ObjectIDFromHex(data.Get(key))
			if err != nil {
				util.LogError("not a valid id", err)
				return false, err
			}
			filter = bson.M{"_id":_id }
		}else{
			value[key] = data.Get(key)
		}
	}
	fmt.Println(filter,bson.M{"$set":value})
	productId, err := c.Collection.UpdateOne(c.Ctx,filter,bson.M{"$set":value})

	if err != nil {
		util.LogError("unable to update product", err)
		return false, err
	}
	if productId.ModifiedCount > 0{
		return  true, nil
	}

	return false, errors.New("data already updated")
}
func (c *Config) Delete(id string) (bool, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		util.LogError("not a valid id", err)
		return false, err
	}
	filter := bson.M{"_id":_id }

	status , err := c.Collection.DeleteOne(c.Ctx, filter)

	if err != nil {
		util.LogError("unable to delete product", err)
		return false, err
	}
	if status.DeletedCount > 0{
		return true, nil
	}
	return false, errors.New("no product to remove")

}
func (c *Config) FindAllProduct(startIndex *int64, limit *int64) (interface{}, error) {

	var product bson.A

	findOpt := &options.FindOptions{
		Limit: limit,
		Skip: startIndex,
	}

	 cur, err := c.Collection.Find(c.Ctx,bson.M{},findOpt)
	 defer cur.Close(c.Ctx)

	if err != nil {
		util.LogError("unable to get product", err)
		return nil, err
	}

	//err = cur.All(c.Ctx,&product)
	for cur.Next(c.Ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {  }
		product = append(product,result)
	}

	if err != nil {
		util.LogError("unable to cur err in product", err)
		return nil, err
	}
	return product, nil

}
func (c *Config) FindSingleProduct(id string) (interface{}, error) {
	 _id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		util.LogError("not a valid id to get product", err)
		return false, err
	}
	filter := bson.M{"_id":_id }
	var result bson.M
	err = c.Collection.FindOne(c.Ctx,filter).Decode(&result)

	if err != nil {
		util.LogError("unable to get product", err)
		return false, err
	}

	return result, nil

}
