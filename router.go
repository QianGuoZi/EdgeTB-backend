package EdgeTB_backend

import "github.com/gin-gonic/gin"

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/EdgeTB")

	//用户 apis
	{
		apiRouter.POST("/login", handler.Login)
	}
}
