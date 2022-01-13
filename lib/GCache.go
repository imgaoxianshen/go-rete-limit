package lib

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type cacheData struct {
	key string
	val interface{}
	expireAt time.Time
}

func newCacheData(key string, val interface{}, expireAt time.Time) *cacheData {
	return &cacheData{key: key, val: val, expireAt: expireAt}
}

type GCacheOption func(g *GCache)
type GCacheOptions []GCacheOption

func (this GCacheOptions) apply(g *GCache) {
	for _, fn := range this{
		fn(g)
	}
}

type GCache struct {
	maxsize int // 限制最大key的数量 0 表示不限制
	elist *list.List
	edata map[string]*list.Element
	lock sync.Mutex
}

func WithMaxSize(size int) GCacheOption{
	return func(g *GCache) {
		if size > 0{
			g.maxsize = size
		}
	}
}

func NewGCache(opt ...GCacheOption) *GCache{
	cache := &GCache{elist: list.New(), edata: make(map[string]*list.Element), maxsize: 10000}
	GCacheOptions(opt).apply(cache)
	cache.clear()
	return cache
}

//获取缓存
func(this *GCache) Get(key string) interface{}{
	this.lock.Lock()
	defer this.lock.Unlock()
	if v, ok := this.edata[key]; ok {
		if time.Now().After(v.Value.(*cacheData).expireAt){ // 过期
			return nil
		}
		this.elist.MoveToFront(v)
		return v.Value.(*cacheData).val
	}
	return nil
}

const NotExpireTTL = time.Hour*24*365 // 不过期的时间
func(this *GCache) Set(key string, newVal interface{}, ttl time.Duration){
	this.lock.Lock()
	defer this.lock.Unlock()
	var setExpire time.Time
	if ttl == 0 {
		setExpire = time.Now().Add(NotExpireTTL)
	}else{
		setExpire = time.Now().Add(ttl)
	}
	newCache := newCacheData(key, newVal, setExpire)
	if v,ok := this.edata[key]; ok{
		v.Value = newCache
		this.elist.MoveToFront(v)
	}else{
		this.edata[key] = this.elist.PushFront(newCache)
		// 判断长度是否溢出 如果末尾淘汰一个缓存
		if this.maxsize > 0 && len(this.edata) > this.maxsize{
			this.removeOldest()
		}
	}
}

//删除最后一个元素
func (this *GCache) removeOldest(){
	back := this.elist.Back()
	if back == nil {
		return
	}
	this.removeItem(back)
}

func(this *GCache) clear() {
	go func() {
		for{
			this.removeExpired()
			time.Sleep(time.Second)
		}
	}()
}

func (this *GCache) removeItem(ele *list.Element) {
	key := ele.Value.(*cacheData).key
	delete(this.edata, key) // 删除map里面的key
	this.elist.Remove(ele)
}

func(this *GCache) removeExpired() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _,v := range this.edata{
		if time.Now().After(v.Value.(*cacheData).expireAt) {
			this.removeItem(v)
		}
	}
}

func(this *GCache) Print() {
	ele := this.elist.Front()
	if ele == nil {
		return
	}

	for {
		fmt.Println(this.Get(ele.Value.(*cacheData).key))
		ele = ele.Next()
		if ele == nil {
			break
		}
	}
}