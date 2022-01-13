package lib

import "github.com/gin-gonic/gin"

func Limiter(cap int64) func(handler gin.HandlerFunc) gin.HandlerFunc{
	limiter := NewBucket(cap,1)
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(context *gin.Context) {
			if limiter.IsAccept() {
				handler(context)
			}else{
				context.AbortWithStatusJSON(429, gin.H{"message":"too many requests"})
			}
		}
	}
}