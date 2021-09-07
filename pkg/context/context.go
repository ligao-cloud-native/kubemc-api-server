package context

import (
	gocontext "context"
)

// ModuleContext is interface for context module management
type ModuleContext interface {
	AddModule(module string)
	AddModuleGroup(module, group string)
	Cleanup(module string)
}

// Context is global context object
type beehiveContext struct {
	moduleContext ModuleContext
	//messageContext MessageContext
	ctx    gocontext.Context
	cancel gocontext.CancelFunc
}
