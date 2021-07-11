package middlewares

import "github.com/gin-gonic/gin"

func MultipartMiddleware() gin.HandlerFunc {
	return func(cxt *gin.Context) {
		cxt.Request.ParseMultipartForm(32 << 20)
	}
}
