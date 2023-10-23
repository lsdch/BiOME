package router

type RouterConfig struct {
	BasePath string
}

var Config = RouterConfig{
	BasePath: "/api/v1",
}
