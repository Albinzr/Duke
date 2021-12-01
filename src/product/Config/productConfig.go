package ProductConfig

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Database       *mongo.Database
	CollectionName string
	Collection     *mongo.Collection
	Ctx            context.Context
}

//material string
//color string
//size int8
//placeholderInfo []PlaceholderDetails
//type PlaceholderDetails struct {
//	width int
//	height int
//	x int
//	y int
//}
