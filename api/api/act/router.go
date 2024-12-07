package act

import "github.com/gin-gonic/gin"

func InitActRouter(r *gin.RouterGroup) {
	r.GET("/act", GinWsHandler)
}
