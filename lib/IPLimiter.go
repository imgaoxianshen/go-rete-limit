package lib

import (
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"time"
)

type LimiterCache struct {
	data sync.Map // key -> ip + 端口 value => bucket
}

var IpCache *LimiterCache
var IpCache2 *GCache

func init() {
	IpCache = &LimiterCache{}
	IpCache2 = NewGCache(WithMaxSize(10000))
}

func IPLimiter(cap, rate int64) func(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(context *gin.Context) {
			ip := context.Request.RemoteAddr
			var limiter *Bucket
			//if v, ok := IpCache.data.Load(ip); ok {
			//	limiter = v.(*Bucket)
			//}else{
			//	limiter = NewBucket(cap, rate)
			//	IpCache.data.Store(ip, limiter)
			//}
			if v := IpCache2.Get(ip); v != nil {
				limiter = v.(*Bucket)
			}else{
				limiter = NewBucket(cap, rate)
				log.Println("from cache")
				IpCache2.Set(ip,limiter, time.Second * 5)
			}

			if limiter.IsAccept() {
				handler(context)
			}else{
				context.AbortWithStatusJSON(429, gin.H{"message": "too many requests"})
			}
		}
	}
}