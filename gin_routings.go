package ezutil

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Group    string
	Versions []RouteVersionConfig
	Handlers []gin.HandlerFunc
}

type RouteVersionConfig struct {
	Version  int
	Groups   []RouteGroupConfig
	Handlers []gin.HandlerFunc
}

type RouteGroupConfig struct {
	Group     string
	Endpoints []EndpointConfig
	Handlers  []gin.HandlerFunc
}

type EndpointConfig struct {
	Method   string
	Endpoint string
	Handlers []gin.HandlerFunc
}

func SetupRoutes(router *gin.Engine, routeConfigs []RouteConfig) {
	if router == nil {
		log.Fatal("Router cannot be nil")
	}

	for _, routeConfig := range routeConfigs {
		routeGroup := router.Group(routeConfig.Group, routeConfig.Handlers...)
		for _, versionConfig := range routeConfig.Versions {
			versionGroup := routeGroup.Group(fmt.Sprintf("/v%d", versionConfig.Version), versionConfig.Handlers...)
			for _, routeGroupConfig := range versionConfig.Groups {
				group := versionGroup.Group(routeGroupConfig.Group, routeGroupConfig.Handlers...)
				for _, endpointConfig := range routeGroupConfig.Endpoints {
					group.Handle(endpointConfig.Method, endpointConfig.Endpoint, endpointConfig.Handlers...)
				}
			}

		}
	}
}
