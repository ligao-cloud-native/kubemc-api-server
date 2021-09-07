package context

import (
	gocontext "context"
	"k8s.io/klog/v2"
	"sync"
)

// define channel type
const (
	MsgCtxTypeChannel = "channel"
)

var (
	// singleton
	context *beehiveContext
	once    sync.Once
)

// InitContext gets global context instance
func InitContext(contextType string) {
	once.Do(func() {
		ctx, cancel := gocontext.WithCancel(gocontext.Background())
		context = &beehiveContext{
			ctx:    ctx,
			cancel: cancel,
		}
		switch contextType {
		case MsgCtxTypeChannel:
			//channelContext := NewChannelContext()
			//context.messageContext = channelContext
			//context.moduleContext = channelContext
		default:
			klog.Fatalf("Do not support context type:%s", contextType)
		}
	})
}

func GetContext() gocontext.Context {
	return context.ctx
}
func Done() <-chan struct{} {
	return context.ctx.Done()
}

// Cancel function
func Cancel() {
	context.cancel()
}
