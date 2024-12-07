package router

import (
	"github.com/gin-gonic/gin"
	"github.com/polaris/codesandbox/api/api/act"
	"net/http"
)

func InitRouter() *gin.Engine {
	// 引入gin框架
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
	})

	root := r.Group("/api")
	{
		// Authorization required and websocket request
		w := root.Group("/")
		{
			act.InitActRouter(w)
		}

		// Authorization required and not websocket request
		//g := root.Group("/")
		//{
		//
		//}
	}

	//// 注册路由
	////http.HandleFunc("/run-code", api.HandleRunCode)
	//http.HandleFunc("/run-code", act.WsHandler)
	return r
}
