package aws

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RouteInfo represents the structure to hold the route details
type RouteInfo struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// GetRoutesHandler is an API that exposes the registered routes of the Gin engine
func GetRoutesHandler(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		routes := router.Routes()
		var routeList []RouteInfo
		for _, r := range routes {
			routeList = append(routeList, RouteInfo{
				Method: r.Method,
				Path:   r.Path,
			})
		}
		c.JSON(http.StatusOK, routeList)
	}
}
