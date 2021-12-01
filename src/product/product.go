package product

import (
	"duke/init/src/product/Config"
	"duke/init/src/product/router"
)

type Config ProductConfig.Config

func (c *Config) Init() {
	routerConfig := (*router.Config)(c)
	routerConfig.Init()
}
