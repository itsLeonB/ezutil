package ezutil

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// RouteConfig represents a top-level route group configuration.
// It defines a route group with optional versioned subgroups and middleware handlers.
// Used by SetupRoutes to create hierarchical route structures.
type RouteConfig struct {
	Group    string
	Versions []RouteVersionConfig
	Handlers []gin.HandlerFunc
}

// RouteVersionConfig represents a versioned route group within a RouteConfig.
// It contains version-specific route groups and middleware handlers.
// Versions are typically represented as integers (e.g., 1 for "/v1").
type RouteVersionConfig struct {
	Version  int
	Groups   []RouteGroupConfig
	Handlers []gin.HandlerFunc
}

// RouteGroupConfig represents a route group within a versioned section.
// It contains individual endpoints and middleware handlers specific to this group.
// Groups help organize related endpoints under a common path prefix.
type RouteGroupConfig struct {
	Group     string
	Endpoints []EndpointConfig
	Handlers  []gin.HandlerFunc
}

// EndpointConfig represents an individual HTTP endpoint.
// It specifies the HTTP method, endpoint path, and handlers for a specific route.
// This is the leaf level of the routing hierarchy.
type EndpointConfig struct {
	Method   string
	Endpoint string
	Handlers []gin.HandlerFunc
}

// SetupRoutes configures Gin routes based on the provided RouteConfig slice.
// It creates a hierarchical route structure with groups, versions, and endpoints.
// The function applies middleware handlers at each level of the hierarchy.
// Panics if router is nil.
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
