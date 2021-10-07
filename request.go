package ginWrapper

import "github.com/gin-gonic/gin"

type GinContext struct {
	*gin.Context
}

func NewGinContext(c *gin.Context) *GinContext {
	return &GinContext{
		Context: c,
	}
}

type GinContextFunc func(ctx *GinContext)

func WrapGinContext(handler GinContextFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(NewGinContext(ctx))
	}
}
