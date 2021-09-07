package ttl

import (
	"sync"
	"time"
)

// token ttl
type TokenTTL struct {
	token    map[string]*item
	ttl      int
	lock     sync.Mutex
	callback func(object string)
}

type item struct {
	value      interface{}
	lastAccess int64
}

func New(maxTTL int, callback func(k string)) (t *TokenTTL) {
	t = &TokenTTL{
		token:    make(map[string]*item),
		ttl:      maxTTL,
		callback: callback,
	}

	// 周期任务，每秒去检测token是否过期
	go func() {
		for now := range time.Tick(time.Second) {
			t.lock.Lock()
			for k, v := range t.token {
				// 删除过期token
				if now.Unix()-v.lastAccess > int64(t.ttl) {
					delete(t.token, k)
					go t.callback(k)
				}
			}
			t.lock.Unlock()
		}
	}()

	return
}

func (t *TokenTTL) Get(k string) (v interface{}) {
	t.lock.Lock()
	if item, ok := t.token[k]; ok {
		v = item.value
		item.lastAccess = time.Now().Unix()
	}
	t.lock.Unlock()
	return

}

func (t *TokenTTL) PUT(k, v string) {
	t.lock.Lock()

	//TODO: 每次调用都需要更新token最后访问时间，如果存在值可能一直不变。
	it, ok := t.token[k]
	if !ok {
		it = &item{value: v}
		t.token[k] = it
	}
	it.lastAccess = time.Now().Unix()

	t.lock.Unlock()
}
