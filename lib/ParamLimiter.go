package lib

import "github.com/gin-gonic/gin"

func ParamLimiter(cap, rate int64,key string) func(handler gin.HandlerFunc) gin.HandlerFunc{
	limiter := NewBucket(cap, rate)
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(context *gin.Context) {
			if context.Query(key) != "" {
				if limiter.IsAccept(){
					handler(context)
				}else{
					context.AbortWithStatusJSON(429, gin.H{"message":"too many rate paramlimit"})
				}
			}else{
				handler(context)
			}
		}
	}
}