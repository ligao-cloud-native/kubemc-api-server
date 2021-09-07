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
				delete(t.token, k)
				// 删除过期token
				if now.Unix()-v.lastAccess > int64(t.ttl) {
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
