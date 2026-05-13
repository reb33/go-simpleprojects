package api

import "3-struct/config"

type Api struct {
	config config.Config
}

func NewApi(config config.Config) *Api{
	return &Api{
		config: config,
	}
}