package main

import (
	"fmt"
	"ratelimit/lib"
	"time"
)

func main() {
	cache := lib.NewGCache(lib.WithMaxSize(3))
	cache.Set("name", "gaofeng", time.Second * 3)
	cache.Set("age", 0, 0)
	for {
		fmt.Printf("name=%v,age=%v",cache.Get("name"), cache.Get("age"))
		time.Sleep(time.Second)
	}
}
