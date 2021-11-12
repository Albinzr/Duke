package cache

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
)

const (
	online = "online-"
	today  = "today-"
)

//Config :- config for redis
type Config struct {
	Host string
	DB int
	Port string
	Password string
	MaxRetries int
	client *redis.Client
}

var ctx = context.Background()

//Init :- init cache
func (c *Config) Init() {
	c.client = redis.NewClient(&redis.Options{
		Addr:      c.Host + ":" + c.Port,
		DB:        c.DB,
		Password: c.Password,
		MaxRetries: c.MaxRetries,
		OnConnect: onConnect,
	})
	c.client.Ping(ctx)
}

func onConnect(ctx context.Context, cn *redis.Conn) error {
	fmt.Println("redis connected")
	return nil
}